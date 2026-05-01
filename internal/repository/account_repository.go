package repository

import (
	"AccountFlow/internal/domain"
	"context"
	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	FindByID(ctx context.Context, accountID uuid.UUID) (*domain.Account, error)
}
