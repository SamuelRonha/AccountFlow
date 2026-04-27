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

func transactionUC(
	txRepo *mocks.MockTransactionRepository,
	accountRepo *mocks.MockAccountRepository,
	opTypeRepo *mocks.MockOperationTypeRepository,
) *usecase.TransactionUseCase {
	return usecase.NewTransactionUseCase(txRepo, accountRepo, opTypeRepo)
}

func stubAccount(id uuid.UUID) *mocks.MockAccountRepository {
	return &mocks.MockAccountRepository{
		FindByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) {
			return &domain.Account{AccountID: id, DocumentNumber: "12345678900"}, nil
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

func stubTxCreate() *mocks.MockTransactionRepository {
	return &mocks.MockTransactionRepository{
		CreateFn: func(_ context.Context, _ *domain.Transaction) error { return nil },
	}
}

// ── CreateTransaction ─────────────────────────────────────────────────────────

func TestCreateTransaction(t *testing.T) {
	accountID := uuid.New()

	// ── sign/amount validation (table-driven) ─────────────────────────────
	signCases := []struct {
		name    string
		opID    int
		opDesc  string
		amount  float64
		wantErr error
		wantAmt float64
	}{
		// Happy path — each op type with correct sign
		{"normal purchase stores -50.0",       1, "Normal Purchase",            -50.0, nil, -50.0},
		{"installments stores -23.5",          2, "Purchase with Installments", -23.5, nil, -23.5},
		{"withdrawal stores -18.7",            3, "Withdrawal",                 -18.7, nil, -18.7},
		{"credit voucher stores +60.0",        4, "Credit Voucher",             60.0,  nil, 60.0},

		// Wrong sign
		{"debit with positive amount rejected", 1, "Normal Purchase",  50.0,  domain.ErrInvalidAmount, 0},
		{"credit with negative amount rejected",4, "Credit Voucher",   -60.0, domain.ErrInvalidAmount, 0},

		// Zero
		{"zero amount rejected", 1, "Normal Purchase", 0, domain.ErrInvalidAmount, 0},
	}

	for _, tc := range signCases {
		t.Run(tc.name, func(t *testing.T) {
			txRepo := stubTxCreate()
			uc := transactionUC(txRepo, stubAccount(accountID), stubOpType(tc.opID, tc.opDesc))

			tx, err := uc.CreateTransaction(context.Background(), accountID, tc.opID, tc.amount)

			if err != tc.wantErr {
				t.Fatalf("got error %v, want %v", err, tc.wantErr)
			}
			if err == nil {
				if tx.Amount != tc.wantAmt {
					t.Errorf("got amount %f, want %f", tx.Amount, tc.wantAmt)
				}
				if txRepo.CreateCalls != 1 {
					t.Errorf("Create called %d times, want 1", txRepo.CreateCalls)
				}
			}
		})
	}

	// ── infrastructure error cases ────────────────────────────────────────

	t.Run("error — account not found", func(t *testing.T) {
		uc := transactionUC(
			&mocks.MockTransactionRepository{},
			&mocks.MockAccountRepository{
				FindByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Account, error) {
					return nil, domain.ErrAccountNotFound
				},
			},
			&mocks.MockOperationTypeRepository{},
		)
		_, err := uc.CreateTransaction(context.Background(), uuid.New(), 1, -50.0)
		if err != domain.ErrAccountNotFound {
			t.Errorf("got %v, want ErrAccountNotFound", err)
		}
	})

	t.Run("error — operation type not found", func(t *testing.T) {
		uc := transactionUC(
			&mocks.MockTransactionRepository{},
			stubAccount(accountID),
			&mocks.MockOperationTypeRepository{
				FindByIDFn: func(_ context.Context, _ int) (*domain.OperationType, error) {
					return nil, domain.ErrOperationTypeNotFound
				},
			},
		)
		_, err := uc.CreateTransaction(context.Background(), accountID, 99, -50.0)
		if err != domain.ErrOperationTypeNotFound {
			t.Errorf("got %v, want ErrOperationTypeNotFound", err)
		}
	})
}
