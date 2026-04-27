package usecase

import (
	"context"

	"AccountFlow/internal/domain"
	"AccountFlow/internal/repository"
	"github.com/google/uuid"
)

type AccountUseCase struct {
	accountRepo repository.AccountRepository
}

func NewAccountUseCase(accountRepo repository.AccountRepository) *AccountUseCase {
	return &AccountUseCase{accountRepo: accountRepo}
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
func (uc *AccountUseCase) GetByID(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	return uc.accountRepo.FindByID(ctx, accountID)
}
