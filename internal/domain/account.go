package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	AccountID      uuid.UUID `json:"account_id"`
	DocumentNumber string    `json:"document_number"`
	Balance        float64   `json:"balance"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewAccount(documentNumber string, balance float64) (*Account, error) {
	documentNumber = strings.TrimSpace(documentNumber)
	if documentNumber == "" {
		return nil, ErrInvalidField
	}
	if balance <= 0 {
		return nil, ErrInvalidBalance
	}
	return &Account{
		AccountID:      uuid.New(),
		DocumentNumber: documentNumber,
		Balance:        balance,
		CreatedAt:      time.Now().UTC(),
	}, nil
}
