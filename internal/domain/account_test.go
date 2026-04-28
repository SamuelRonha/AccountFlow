package domain_test

import (
	"testing"

	"AccountFlow/internal/domain"
	"github.com/google/uuid"
)

// ── Account ──────────────────────────────────────────────────────────────────

func TestNewAccount(t *testing.T) {
	t.Run("valid document creates account", func(t *testing.T) {
		acc, err := domain.NewAccount("12345678900")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.DocumentNumber != "12345678900" {
			t.Errorf("got document %q, want %q", acc.DocumentNumber, "12345678900")
		}
		if acc.AccountID == uuid.Nil {
			t.Error("AccountID must not be nil")
		}
	})

	invalidCases := []struct {
		name     string
		document string
	}{
		{"empty string", ""},
		{"only whitespace", "   "},
	}
	for _, tc := range invalidCases {
		t.Run("rejects "+tc.name, func(t *testing.T) {
			_, err := domain.NewAccount(tc.document)
			if err != domain.ErrInvalidField {
				t.Errorf("got %v, want ErrInvalidField", err)
			}
		})
	}
}

// ── Transaction ───────────────────────────────────────────────────────────────

func TestNewTransaction_SignValidation(t *testing.T) {
	// The caller always sends a POSITIVE amount.
	// The domain applies the sign automatically based on operation type:
	//   Debit  (type 1, 2, 3) → stored as negative  (200 → -200)
	//   Credit (type 4)       → stored as positive   (60  → +60)
	cases := []struct {
		name       string
		opID       int
		amount     float64
		wantErr    error
		wantStored float64
	}{
		// ── happy path — caller sends positive ─────────────────────────────
		{"normal purchase 50 stored as -50", 1, 50.0, nil, -50.0},
		{"installments 23.5 stored as -23.5", 2, 23.5, nil, -23.5},
		{"withdrawal 18.7 stored as -18.7", 3, 18.7, nil, -18.7},
		{"credit voucher 60 stored as +60", 4, 60.0, nil, 60.0},

		// ── invalid input ───────────────────────────────────────────────────
		{"zero amount rejected", 1, 0, domain.ErrInvalidAmount, 0},
		{"negative amount rejected", 1, -50.0, domain.ErrInvalidAmount, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opType := &domain.OperationType{OperationTypeID: tc.opID}
			tx, err := domain.NewTransaction(uuid.New(), opType, tc.amount)

			if err != tc.wantErr {
				t.Fatalf("got error %v, want %v", err, tc.wantErr)
			}
			if err == nil && tx.Amount != tc.wantStored {
				t.Errorf("got stored amount %f, want %f", tx.Amount, tc.wantStored)
			}
		})
	}
}

// ── OperationType ────────────────────────────────────────────────────────────

func TestOperationType_IsDebit(t *testing.T) {
	cases := []struct {
		opID    int
		desc    string
		isDebit bool
	}{
		{1, "Normal Purchase", true},
		{2, "Purchase with Installments", true},
		{3, "Withdrawal", true},
		{4, "Credit Voucher", false},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			op := &domain.OperationType{OperationTypeID: tc.opID, Description: tc.desc}
			if op.IsDebit() != tc.isDebit {
				t.Errorf("IsDebit() = %v, want %v", op.IsDebit(), tc.isDebit)
			}
		})
	}
}
