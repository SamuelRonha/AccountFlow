package usecase

import (
	"context"

	"AccountFlow/internal/domain"
	"AccountFlow/internal/repository"
	"github.com/google/uuid"
)

type AccountUseCase struct {
	accountRepo repository.AccountRepository
	txRepo      repository.TransactionRepository
}

func NewAccountUseCase(accountRepo repository.AccountRepository, txRepo repository.TransactionRepository) *AccountUseCase {
	return &AccountUseCase{accountRepo: accountRepo, txRepo: txRepo}
}

// CreateAccount opens a new account for the given document number.
// Returns ErrDocumentAlreadyUsed if the document is already registered.
func (uc *AccountUseCase) CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	account, err := domain.NewAccount(documentNumber)
	if err != nil {
		return nil, err
	}

	if err := uc.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

// GetByID retrieves an account by its ID.
func (uc *AccountUseCase) GetByID(ctx context.Context, accountID uuid.UUID) (*domain.AccountResponse, error) {
	acc, err := uc.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	balance, err := uc.txRepo.GetBalance(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return &domain.AccountResponse{
		AccountID:      acc.AccountID,
		DocumentNumber: acc.DocumentNumber,
		CreatedAt:      acc.CreatedAt,
		Balance:        balance,
	}, nil
}
