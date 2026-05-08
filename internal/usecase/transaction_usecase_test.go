package usecase_test

import (
	"context"
	"database/sql"
	"testing"

	"AccountFlow/internal/domain"
	"AccountFlow/internal/usecase"
	"AccountFlow/internal/usecase/mocks"
	"github.com/google/uuid"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func transactionUC(
	txRepo *mocks.MockTransactionRepository,
	accountRepo *mocks.MockAccountRepository,
	opTypeRepo *mocks.MockOperationTypeRepository,
) *usecase.TransactionUseCase {
	return usecase.NewTransactionUseCase(txRepo, accountRepo, opTypeRepo)
}

// stubAccount returns an account mock whose UpdateBalance always succeeds.
func stubAccount() *mocks.MockAccountRepository {
	return &mocks.MockAccountRepository{
		UpdateBalanceFn: func(_ context.Context, _ *sql.Tx, _ uuid.UUID, _ float64) error {
			return nil
		},
	}
}

func stubOpType(id int, desc string) *mocks.MockOperationTypeRepository {
	return &mocks.MockOperationTypeRepository{
		FindByIDFn: func(_ context.Context, _ int) (*domain.OperationType, error) {
			return &domain.OperationType{OperationTypeID: id, Description: desc}, nil
		},
	}
}

func stubTxRepo() *mocks.MockTransactionRepository {
	return &mocks.MockTransactionRepository{
		// BeginTxFn intentionally unset — mock returns nil tx (no-op commit/rollback in usecase).
		CreateTxFn: func(_ context.Context, _ *sql.Tx, _ *domain.Transaction) error { return nil },
	}
}

// ── CreateTransaction ─────────────────────────────────────────────────────────

func TestCreateTransaction(t *testing.T) {
	accountID := uuid.New()

	// ── sign/amount validation (table-driven) ─────────────────────────────
	signCases := []struct {
		name         string
		opID         int
		opDesc       string
		updateBalErr error   // what UpdateBalance mock returns
		amount       float64 // caller always sends positive
		wantErr      error
		wantAmt      float64 // what gets stored (sign applied by domain)
	}{
		// Happy path — caller sends positive, domain applies sign
		{"normal purchase 50 stored as -50", 1, "Normal Purchase", nil, 50.0, nil, -50.0},
		{"installments 23.5 stored as -23.5", 2, "Purchase with Installments", nil, 23.5, nil, -23.5},
		{"withdrawal 18.7 stored as -18.7", 3, "Withdrawal", nil, 18.7, nil, -18.7},
		{"credit voucher 60 stored as +60", 4, "Credit Voucher", nil, 60.0, nil, 60.0},

		// Invalid input (rejected before hitting the DB)
		{"zero amount rejected", 1, "Normal Purchase", nil, 0, domain.ErrInvalidAmount, 0},
		{"negative amount rejected", 1, "Normal Purchase", nil, -50.0, domain.ErrInvalidAmount, 0},

		// DB-enforced balance guard
		{"insufficient balance rejected", 1, "Normal Purchase", domain.ErrInsufficientBalance, 50.0, domain.ErrInsufficientBalance, 0},
	}

	for _, tc := range signCases {
		t.Run(tc.name, func(t *testing.T) {
			txRepo := stubTxRepo()
			accountRepo := &mocks.MockAccountRepository{
				UpdateBalanceFn: func(_ context.Context, _ *sql.Tx, _ uuid.UUID, _ float64) error {
					return tc.updateBalErr
				},
			}
			uc := transactionUC(txRepo, accountRepo, stubOpType(tc.opID, tc.opDesc))

			tx, err := uc.CreateTransaction(context.Background(), accountID, tc.opID, tc.amount)

			if err != tc.wantErr {
				t.Fatalf("got error %v, want %v", err, tc.wantErr)
			}
			if err == nil {
				if tx.Amount != tc.wantAmt {
					t.Errorf("got amount %f, want %f", tx.Amount, tc.wantAmt)
				}
				if txRepo.CreateTxCalls != 1 {
					t.Errorf("CreateTx called %d times, want 1", txRepo.CreateTxCalls)
				}
			}
		})
	}

	// ── infrastructure error cases ────────────────────────────────────────

	t.Run("error — account not found", func(t *testing.T) {
		uc := transactionUC(
			stubTxRepo(),
			&mocks.MockAccountRepository{
				UpdateBalanceFn: func(_ context.Context, _ *sql.Tx, _ uuid.UUID, _ float64) error {
					return domain.ErrAccountNotFound
				},
			},
			stubOpType(1, "Normal Purchase"),
		)
		_, err := uc.CreateTransaction(context.Background(), uuid.New(), 1, 50.0)
		if err != domain.ErrAccountNotFound {
			t.Errorf("got %v, want ErrAccountNotFound", err)
		}
	})

	t.Run("error — operation type not found", func(t *testing.T) {
		uc := transactionUC(
			stubTxRepo(),
			stubAccount(),
			&mocks.MockOperationTypeRepository{
				FindByIDFn: func(_ context.Context, _ int) (*domain.OperationType, error) {
					return nil, domain.ErrOperationTypeNotFound
				},
			},
		)
		_, err := uc.CreateTransaction(context.Background(), accountID, 99, 50.0)
		if err != domain.ErrOperationTypeNotFound {
			t.Errorf("got %v, want ErrOperationTypeNotFound", err)
		}
	})
}
