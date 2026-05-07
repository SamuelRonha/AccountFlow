package domain

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

type Transfer struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	IdTransfer    string    `json:"id_transfer"`
}

func NewTransference(from, to uuid.UUID, amount float64, transferID string) (*Transfer, error) {
	if from == uuid.Nil || to == uuid.Nil {
		return nil, ErrInvalidField
	}

	transferID = strings.TrimSpace(transferID)
	if transferID == "" {
		return nil, ErrInvalidField
	}

	if from == to {
		return nil, errors.New("from and to accounts must be different")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	return &Transfer{
		FromAccountID: from,
		ToAccountID:   to,
		Amount:        amount,
		IdTransfer:    transferID,
	}, nil
}
