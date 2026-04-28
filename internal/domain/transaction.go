package domain

import (
	"time"

	"github.com/google/uuid"
)

// Transaction records a financial movement on an account.
//
// The Amount field uses the following sign convention (enforced by this package):
//   - Negative → debit operations (Normal Purchase, Installments, Withdrawal)
//   - Positive → credit operations (Credit Voucher)
type Transaction struct {
	TransactionID   uuid.UUID `json:"transaction_id"`
	AccountID       uuid.UUID `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}

// NewTransaction creates a transaction from a caller-supplied positive amount.
// The domain applies the correct sign automatically based on operation type:
//   - Debit operations  (type 1, 2, 3) → amount is stored as negative (e.g. 50.0 → -50.0)
//   - Credit operations (type 4)       → amount is stored as positive (e.g. 60.0 → +60.0)
//
// The caller must always send a positive value. Returns ErrInvalidAmount if
// the amount is zero or negative.
func NewTransaction(accountID uuid.UUID, opType *OperationType, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if opType.IsDebit() {
		amount = -amount
	}

	return &Transaction{
		TransactionID:   uuid.New(),
		AccountID:       accountID,
		OperationTypeID: opType.OperationTypeID,
		Amount:          amount,
		EventDate:       time.Now().UTC(),
	}, nil
}
