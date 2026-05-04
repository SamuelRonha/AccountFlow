package postgres

import (
	"context"
	"database/sql"
	"errors"

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
	ts, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	isCommitted := false

	defer func() {
		if !isCommitted {
			_ = ts.Rollback()
		}
	}()
	var exists int
	queryLock := `SELECT 1 FROM accounts WHERE account_id = $1 FOR UPDATE`

	err = ts.QueryRowContext(ctx, queryLock, tx.AccountID).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrAccountNotFound
		}
		return err
	}

	_, err = ts.ExecContext(ctx,
		`INSERT INTO transactions (transaction_id, account_id, operation_type_id, amount, event_date)
         VALUES ($1, $2, $3, $4, $5)`,
		tx.TransactionID, tx.AccountID, tx.OperationTypeID, tx.Amount, tx.EventDate,
	)
	if err != nil {
		return err
	}

	if err = ts.Commit(); err != nil {
		return err
	}
	isCommitted = true

	return nil
}

func (r *TransactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error) {
	query := `SELECT transaction_id, account_id, operation_type_id, amount, event_date
	          FROM transactions WHERE account_id = $1 ORDER BY event_date DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			//log but not subscribe the main error :>
		}
	}()

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
