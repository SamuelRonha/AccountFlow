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
	// The caller is responsible for the sign; the domain validates it matches
	// the operation type.
	//
	// Debit operations (type 1, 2, 3) → amount must be NEGATIVE
	// Credit operations (type 4)      → amount must be POSITIVE
	cases := []struct {
		name    string
		opID    int
		amount  float64
		wantErr error
	}{
		// ── happy path ─────────────────────────────────────────────────────
		{"normal purchase with negative amount",       1, -50.0, nil},
		{"installments with negative amount",          2, -23.5, nil},
		{"withdrawal with negative amount",            3, -18.7, nil},
		{"credit voucher with positive amount",        4, 60.0,  nil},

		// ── wrong sign ─────────────────────────────────────────────────────
		{"normal purchase with positive amount",       1, 50.0,  domain.ErrInvalidAmount},
		{"installments with positive amount",          2, 23.5,  domain.ErrInvalidAmount},
		{"withdrawal with positive amount",            3, 18.7,  domain.ErrInvalidAmount},
		{"credit voucher with negative amount",        4, -60.0, domain.ErrInvalidAmount},

		// ── zero ───────────────────────────────────────────────────────────
		{"zero amount on debit op",  1, 0, domain.ErrInvalidAmount},
		{"zero amount on credit op", 4, 0, domain.ErrInvalidAmount},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opType := &domain.OperationType{OperationTypeID: tc.opID}
			tx, err := domain.NewTransaction(uuid.New(), opType, tc.amount)

			if err != tc.wantErr {
				t.Fatalf("got error %v, want %v", err, tc.wantErr)
			}
			if err == nil && tx.Amount != tc.amount {
				t.Errorf("got amount %f, want %f", tx.Amount, tc.amount)
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
		{1, "Normal Purchase",            true},
		{2, "Purchase with Installments", true},
		{3, "Withdrawal",                 true},
		{4, "Credit Voucher",             false},
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
