package repository

import (
	"context"

	"AccountFlow/internal/domain"
	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error)
}
