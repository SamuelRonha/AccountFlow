package postgres

import (
	"context"
	"database/sql"

	"AccountFlow/internal/domain"
	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `INSERT INTO transactions (transaction_id, account_id, operation_type_id, amount, event_date)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		tx.TransactionID, tx.AccountID, tx.OperationTypeID, tx.Amount, tx.EventDate,
	)
	return err
}

func (r *TransactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error) {
	query := `SELECT transaction_id, account_id, operation_type_id, amount, event_date
	          FROM transactions WHERE account_id = $1 ORDER BY event_date DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		if err := rows.Scan(&tx.TransactionID, &tx.AccountID, &tx.OperationTypeID, &tx.Amount, &tx.EventDate); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return txs, rows.Err()
}
