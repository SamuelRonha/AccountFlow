package repository

import (
	"context"

	"AccountFlow/internal/domain"
)

type OperationTypeRepository interface {
	FindByID(ctx context.Context, id int) (*domain.OperationType, error)
}

