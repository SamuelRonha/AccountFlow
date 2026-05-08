package repository

import (
	"AccountFlow/internal/domain"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	// CreateTx inserts a transaction record within an existing database transaction.
	CreateTx(ctx context.Context, dbTx *sql.Tx, tx *domain.Transaction) error
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
}
