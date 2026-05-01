package postgres

import (
	"context"
	"database/sql"
	"errors"

	"AccountFlow/internal/domain"
)

type OperationTypeRepository struct {
	db *sql.DB
}

func NewOperationTypeRepository(db *sql.DB) *OperationTypeRepository {
	return &OperationTypeRepository{db: db}
}

func (r *OperationTypeRepository) FindByID(ctx context.Context, id int) (*domain.OperationType, error) {
	query := `SELECT operation_type_id, description FROM operation_types WHERE operation_type_id = $1`
	var op domain.OperationType
	err := r.db.QueryRowContext(ctx, query, id).Scan(&op.OperationTypeID, &op.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrOperationTypeNotFound
	}
	if err != nil {
		return nil, err
	}
	return &op, nil
}
