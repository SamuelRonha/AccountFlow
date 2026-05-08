package mocks

import (
	"context"
	"database/sql"
	"fmt"

	"AccountFlow/internal/domain"
	"github.com/google/uuid"
)

// ── AccountRepository ────────────────────────────────────────────────────────

// MockAccountRepository is a test double for repository.AccountRepository.
// Set only the Fn fields your test needs — any unset field will panic with a
// descriptive message if accidentally called, catching missing setup early.
type MockAccountRepository struct {
	CreateFn        func(ctx context.Context, account *domain.Account) error
	FindByIDFn      func(ctx context.Context, accountID uuid.UUID) (*domain.Account, error)
	UpdateBalanceFn func(ctx context.Context, tx *sql.Tx, accountID uuid.UUID, delta float64) error

	// Call counters — inspect these in your tests when needed.
	CreateCalls        int
	FindByIDCalls      int
	UpdateBalanceCalls int
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

func (m *MockAccountRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, accountID uuid.UUID, delta float64) error {
	m.UpdateBalanceCalls++
	if m.UpdateBalanceFn == nil {
		panic(fmt.Sprintf("MockAccountRepository.UpdateBalance called but UpdateBalanceFn is not set (call #%d)", m.UpdateBalanceCalls))
	}
	return m.UpdateBalanceFn(ctx, tx, accountID, delta)
}

// ── TransactionRepository ────────────────────────────────────────────────────

// MockTransactionRepository is a test double for repository.TransactionRepository.
type MockTransactionRepository struct {
	BeginTxFn         func(ctx context.Context) (*sql.Tx, error)
	CreateFn          func(ctx context.Context, tx *domain.Transaction) error
	CreateTxFn        func(ctx context.Context, dbTx *sql.Tx, tx *domain.Transaction) error
	FindByAccountIDFn func(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error)

	BeginTxCalls         int
	CreateCalls          int
	CreateTxCalls        int
	FindByAccountIDCalls int
}

func (m *MockTransactionRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	m.BeginTxCalls++
	if m.BeginTxFn == nil {
		// Return nil so the usecase defers a no-op rollback/commit; fine for unit tests.
		return nil, nil
	}
	return m.BeginTxFn(ctx)
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	m.CreateCalls++
	if m.CreateFn == nil {
		panic(fmt.Sprintf("MockTransactionRepository.Create called but CreateFn is not set (call #%d)", m.CreateCalls))
	}
	return m.CreateFn(ctx, tx)
}

func (m *MockTransactionRepository) CreateTx(ctx context.Context, dbTx *sql.Tx, tx *domain.Transaction) error {
	m.CreateTxCalls++
	if m.CreateTxFn == nil {
		panic(fmt.Sprintf("MockTransactionRepository.CreateTx called but CreateTxFn is not set (call #%d)", m.CreateTxCalls))
	}
	return m.CreateTxFn(ctx, dbTx, tx)
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
