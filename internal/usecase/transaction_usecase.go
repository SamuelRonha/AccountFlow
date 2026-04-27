package usecase

import (
	"context"

	"AccountFlow/internal/domain"
	"AccountFlow/internal/repository"
	"github.com/google/uuid"
)

type TransactionUseCase struct {
	txRepo      repository.TransactionRepository
	accountRepo repository.AccountRepository
	opTypeRepo  repository.OperationTypeRepository
}

func NewTransactionUseCase(
	txRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	opTypeRepo repository.OperationTypeRepository,
) *TransactionUseCase {
	return &TransactionUseCase{
		txRepo:      txRepo,
		accountRepo: accountRepo,
		opTypeRepo:  opTypeRepo,
	}
}

// CreateTransaction validates the account and operation type, then persists the transaction.
// The amount sign is enforced by the operation type: debits become negative, credits positive.
func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, accountID uuid.UUID, operationTypeID int, amount float64) (*domain.Transaction, error) {
	if _, err := uc.accountRepo.FindByID(ctx, accountID); err != nil {
		return nil, err
	}

	opType, err := uc.opTypeRepo.FindByID(ctx, operationTypeID)
	if err != nil {
		return nil, err
	}

	tx, err := domain.NewTransaction(accountID, opType, amount)
	if err != nil {
		return nil, err
	}

	if err := uc.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}
