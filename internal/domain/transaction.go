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

// NewTransaction creates a transaction and validates that the amount sign is
// consistent with the operation type:
//   - Debit operations  (type 1, 2, 3) require a negative amount  (e.g. -50.0)
//   - Credit operations (type 4)       require a positive amount  (e.g.  60.0)
//
// Returns ErrInvalidAmount if the amount is zero or the sign contradicts the
// operation type.
func NewTransaction(accountID uuid.UUID, opType *OperationType, amount float64) (*Transaction, error) {
	if amount == 0 {
		return nil, ErrInvalidAmount
	}
	if opType.IsDebit() && amount > 0 {
		return nil, ErrInvalidAmount // debit operations must be negative
	}
	if !opType.IsDebit() && amount < 0 {
		return nil, ErrInvalidAmount // credit operations must be positive
	}

	return &Transaction{
		TransactionID:   uuid.New(),
		AccountID:       accountID,
		OperationTypeID: opType.OperationTypeID,
		Amount:          amount,
		EventDate:       time.Now().UTC(),
	}, nil
}
