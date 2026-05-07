package repository

import (
	"context"
	"database/sql"

	"AccountFlow/internal/domain"
	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	CreateWithTx(ctx context.Context, sqlTx *sql.Tx, tx *domain.Transaction) error
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	GetBalanceTx(ctx context.Context, tx *sql.Tx, accountID uuid.UUID) (float64, error)
	GetBalance(ctx context.Context, accountID uuid.UUID) (float64, error)
	FindTransactionExistence(ctx context.Context, tx *sql.Tx, TransferID string, AccountID uuid.UUID) bool
}
