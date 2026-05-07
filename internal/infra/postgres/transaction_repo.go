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
	query := `INSERT INTO transactions (transaction_id, account_id, operation_type_id, amount, event_date, transfer_id)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		tx.TransactionID, tx.AccountID, tx.OperationTypeID, tx.Amount, tx.EventDate, tx.IdTransfer,
	)
	return err
}

func (r *TransactionRepository) CreateWithTx(ctx context.Context, sqlTx *sql.Tx, tx *domain.Transaction) error {
	query := `INSERT INTO transactions (transaction_id, account_id, operation_type_id, amount, event_date, transfer_id)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := sqlTx.ExecContext(ctx, query,
		tx.TransactionID, tx.AccountID, tx.OperationTypeID, tx.Amount, tx.EventDate, tx.IdTransfer,
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

func (r *TransactionRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *TransactionRepository) GetBalanceTx(ctx context.Context, tx *sql.Tx, accountID uuid.UUID) (float64, error) {
	var balance float64
	err := tx.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE account_id = $1`, accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *TransactionRepository) GetBalance(ctx context.Context, accountID uuid.UUID) (float64, error) {
	var balance float64
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE account_id = $1`, accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *TransactionRepository) FindTransactionExistence(ctx context.Context, tx *sql.Tx, TransferID string, AccountID uuid.UUID) bool {
	query := `SELECT EXISTS (
    				SELECT 1 FROM transactions WHERE transfer_id = $1 and account_id = $2
			  )`

	var exists bool
	err := tx.QueryRowContext(ctx, query, TransferID, AccountID).Scan(&exists)

	if err != nil {
		return true
	}
	return exists
}
