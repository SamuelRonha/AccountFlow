package usecase

import (
	"context"
	"database/sql"

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

// CreateTransaction validates the account and operation type, applies the transaction,
// and atomically updates the account balance inside a single DB transaction.
//
// Balance integrity is enforced by the DB: the UPDATE only applies when
// balance + delta >= 0, eliminating race conditions without SELECT … FOR UPDATE.
func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, accountID uuid.UUID, operationTypeID int, amount float64) (*domain.Transaction, error) {
	dbTx, err := uc.txRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	committed := false
	defer func() {
		if !committed && dbTx != nil {
			_ = dbTx.Rollback()
		}
	}()

	opType, err := uc.opTypeRepo.FindByID(ctx, operationTypeID)
	if err != nil {
		return nil, err
	}

	transaction, err := domain.NewTransaction(accountID, opType, amount)
	if err != nil {
		return nil, err
	}

	// Atomic delta update — the DB rejects the update when balance would drop below 0.
	// Returns ErrAccountNotFound or ErrInsufficientBalance on 0 rows affected.
	if err := uc.accountRepo.UpdateBalance(ctx, dbTx, accountID, transaction.Amount); err != nil {
		return nil, err
	}

	if err := uc.txRepo.CreateTx(ctx, dbTx, transaction); err != nil {
		return nil, err
	}

	if err := commitTx(dbTx); err != nil {
		return nil, err
	}
	committed = true

	return transaction, nil
}

// commitTx commits the transaction when a real *sql.Tx is provided.
// In unit tests the mock may supply nil, which is a no-op.
func commitTx(tx *sql.Tx) error {
	if tx == nil {
		return nil
	}
	return tx.Commit()
}
