package usecase_test

import (
	"context"
	"testing"

	"AccountFlow/internal/domain"
	"AccountFlow/internal/usecase"
	"AccountFlow/internal/usecase/mocks"
	"github.com/google/uuid"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func accountUC(repo *mocks.MockAccountRepository) *usecase.AccountUseCase {
	return usecase.NewAccountUseCase(repo)
}

// ── CreateAccount ─────────────────────────────────────────────────────────────

func TestCreateAccount(t *testing.T) {
	t.Run("success — valid document", func(t *testing.T) {
		repo := &mocks.MockAccountRepository{
			CreateFn: func(_ context.Context, _ *domain.Account) error { return nil },
		}
		acc, err := accountUC(repo).CreateAccount(context.Background(), "12345678900")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.DocumentNumber != "12345678900" {
			t.Errorf("got %q, want %q", acc.DocumentNumber, "12345678900")
		}
		if acc.AccountID == uuid.Nil {
			t.Error("AccountID must not be nil")
		}
		if repo.CreateCalls != 1 {
			t.Errorf("Create called %d times, want 1", repo.CreateCalls)
		}
	})

	t.Run("error — empty document (rejected before hitting repo)", func(t *testing.T) {
		repo := &mocks.MockAccountRepository{} // CreateFn intentionally unset
		_, err := accountUC(repo).CreateAccount(context.Background(), "")
		if err != domain.ErrInvalidField {
			t.Errorf("got %v, want ErrInvalidField", err)
		}
		if repo.CreateCalls != 0 {
			t.Errorf("Create should not be called, but was called %d times", repo.CreateCalls)
		}
	})

	t.Run("error — document already registered", func(t *testing.T) {
		repo := &mocks.MockAccountRepository{
			CreateFn: func(_ context.Context, _ *domain.Account) error {
				return domain.ErrDocumentAlreadyUsed
			},
		}
		_, err := accountUC(repo).CreateAccount(context.Background(), "12345678900")
		if err != domain.ErrDocumentAlreadyUsed {
			t.Errorf("got %v, want ErrDocumentAlreadyUsed", err)
		}
	})
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestGetByID(t *testing.T) {
	accountID := uuid.New()

	t.Run("success — existing account returned", func(t *testing.T) {
		stored := &domain.Account{AccountID: accountID, DocumentNumber: "11122233344"}
		repo := &mocks.MockAccountRepository{
			FindByIDFn: func(_ context.Context, id uuid.UUID) (*domain.Account, error) {
				return stored, nil
			},
		}
		acc, err := accountUC(repo).GetByID(context.Background(), accountID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.AccountID != accountID {
			t.Errorf("got %s, want %s", acc.AccountID, accountID)
		}
	})

	t.Run("error — account not found", func(t *testing.T) {
		repo := &mocks.MockAccountRepository{
			FindByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) {
				return nil, domain.ErrAccountNotFound
			},
		}
		_, err := accountUC(repo).GetByID(context.Background(), uuid.New())
		if err != domain.ErrAccountNotFound {
			t.Errorf("got %v, want ErrAccountNotFound", err)
		}
	})
}
