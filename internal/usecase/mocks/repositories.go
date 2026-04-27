package mocks

import (
	"context"
	"fmt"

	"AccountFlow/internal/domain"
	"github.com/google/uuid"
)

// ── AccountRepository ────────────────────────────────────────────────────────

// MockAccountRepository is a test double for repository.AccountRepository.
// Set only the Fn fields your test needs — any unset field will panic with a
// descriptive message if accidentally called, catching missing setup early.
type MockAccountRepository struct {
	CreateFn   func(ctx context.Context, account *domain.Account) error
	FindByIDFn func(ctx context.Context, accountID uuid.UUID) (*domain.Account, error)

	// Call counters — inspect these in your tests when needed.
	CreateCalls   int
	FindByIDCalls int
}

func (m *MockAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	m.CreateCalls++
	if m.CreateFn == nil {
		panic(fmt.Sprintf("MockAccountRepository.Create called but CreateFn is not set (call #%d)", m.CreateCalls))
	}
	return m.CreateFn(ctx, account)
}

func (m *MockAccountRepository) FindByID(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	m.FindByIDCalls++
	if m.FindByIDFn == nil {
		panic(fmt.Sprintf("MockAccountRepository.FindByID called but FindByIDFn is not set (call #%d)", m.FindByIDCalls))
	}
	return m.FindByIDFn(ctx, accountID)
}

// ── TransactionRepository ────────────────────────────────────────────────────

// MockTransactionRepository is a test double for repository.TransactionRepository.
type MockTransactionRepository struct {
	CreateFn          func(ctx context.Context, tx *domain.Transaction) error
	FindByAccountIDFn func(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error)

	CreateCalls          int
	FindByAccountIDCalls int
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	m.CreateCalls++
	if m.CreateFn == nil {
		panic(fmt.Sprintf("MockTransactionRepository.Create called but CreateFn is not set (call #%d)", m.CreateCalls))
	}
	return m.CreateFn(ctx, tx)
}

func (m *MockTransactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error) {
	m.FindByAccountIDCalls++
	if m.FindByAccountIDFn == nil {
		panic(fmt.Sprintf("MockTransactionRepository.FindByAccountID called but FindByAccountIDFn is not set (call #%d)", m.FindByAccountIDCalls))
	}
	return m.FindByAccountIDFn(ctx, accountID)
}

// ── OperationTypeRepository ──────────────────────────────────────────────────

// MockOperationTypeRepository is a test double for repository.OperationTypeRepository.
type MockOperationTypeRepository struct {
	FindByIDFn func(ctx context.Context, id int) (*domain.OperationType, error)

	FindByIDCalls int
}

func (m *MockOperationTypeRepository) FindByID(ctx context.Context, id int) (*domain.OperationType, error) {
	m.FindByIDCalls++
	if m.FindByIDFn == nil {
		panic(fmt.Sprintf("MockOperationTypeRepository.FindByID called but FindByIDFn is not set (call #%d)", m.FindByIDCalls))
	}
	return m.FindByIDFn(ctx, id)
}
