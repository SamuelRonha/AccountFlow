package domain

// OperationType represents a financial operation category.
// IDs are pre-seeded in the database:
//
//	1 - Normal Purchase       (debit  → negative amount)
//	2 - Purchase in Installments (debit → negative amount)
//	3 - Withdrawal            (debit  → negative amount)
//	4 - Credit Voucher        (credit → positive amount)
type OperationType struct {
	OperationTypeID int    `json:"operation_type_id"`
	Description     string `json:"description"`
}

// IsDebit returns true when the operation should produce a negative amount.
// All types except Credit Voucher (4) are debit operations.
func (o *OperationType) IsDebit() bool {
	return o.OperationTypeID != 4
}

