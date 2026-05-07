package usecase

import (
	"AccountFlow/internal/domain"
	"AccountFlow/internal/repository"
	"context"
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

	tx, err := domain.NewTransaction(accountID, opType, amount, nil)
	if err != nil {
		return nil, err
	}

	if err := uc.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

const DEBIT = 1
const CREDIT = 4

func (uc *TransactionUseCase) TransferAmount(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount float64, transferID string) error {

	tx, err := uc.txRepo.BeginTx(ctx)
	if err != nil {
		return err
	}

	commited := false
	defer func() {
		if !commited {
			_ = tx.Rollback()
		}
	}()

	//order 4 avoid deadlock
	if fromAccountID.String() > toAccountID.String() {
		fromAccountID, toAccountID = toAccountID, fromAccountID
	}

	transferID = transferID + fromAccountID.String() + toAccountID.String()

	existsFrom := uc.accountRepo.FindByIDForUpdate(ctx, tx, fromAccountID)
	if !existsFrom {
		return domain.ErrAccountNotFound
	}

	existsTo := uc.accountRepo.FindByIDForUpdate(ctx, tx, toAccountID)
	if !existsTo {
		return domain.ErrToAccountNotFound
	}

	fromBalance, err := uc.txRepo.GetBalanceTx(ctx, tx, fromAccountID)

	if fromBalance < amount {
		return domain.ErrInsufficientFunds
	}

	exists := uc.txRepo.FindTransactionExistence(ctx, tx, transferID, fromAccountID)

	if exists {
		return nil
	}

	//debit transaction
	op := domain.NewOperationType(DEBIT, "Debit")
	debitTx, err := domain.NewTransaction(fromAccountID, op, amount, &transferID)
	if err != nil {
		return err
	}
	if err = uc.txRepo.CreateWithTx(ctx, tx, debitTx); err != nil {
		return err
	}
	//credit transaction
	op = domain.NewOperationType(CREDIT, "Credit")
	creditTx, err := domain.NewTransaction(toAccountID, op, amount, &transferID)
	if err != nil {
		return err
	}
	if err = uc.txRepo.CreateWithTx(ctx, tx, creditTx); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	commited = true

	return nil

}
