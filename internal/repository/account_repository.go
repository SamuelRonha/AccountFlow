package repository

import (
	"AccountFlow/internal/domain"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	FindByID(ctx context.Context, accountID uuid.UUID) (*domain.Account, error)
	FindByIDForUpdate(ctx context.Context, tx *sql.Tx, accountID uuid.UUID) bool
}
