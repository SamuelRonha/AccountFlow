package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	AccountID      uuid.UUID `json:"account_id"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewAccount(documentNumber string) (*Account, error) {
	documentNumber = strings.TrimSpace(documentNumber)
	if documentNumber == "" {
		return nil, ErrInvalidField
	}
	return &Account{
		AccountID:      uuid.New(),
		DocumentNumber: documentNumber,
		CreatedAt:      time.Now().UTC(),
	}, nil
}
