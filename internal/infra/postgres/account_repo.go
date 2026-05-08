package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq/pqerror"

	"AccountFlow/internal/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `INSERT INTO accounts (account_id, document_number, balance, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, account.AccountID, account.DocumentNumber, account.Balance, account.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation {
			return domain.ErrDocumentAlreadyUsed
		}
		return err
	}
	return nil
}

func (r *AccountRepository) FindByID(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	query := `SELECT account_id, document_number, coalesce(balance, 0) as balance, created_at FROM accounts WHERE account_id = $1`
	var a domain.Account
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(&a.AccountID, &a.DocumentNumber, &a.Balance, &a.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// UpdateBalance atomically applies delta (positive credit, negative debit) to the account balance.
// The update is only applied when the resulting balance would stay >= 0, so no separate
// locking or in-application balance check is needed.
// Returns ErrAccountNotFound / ErrInsufficientBalance when 0 rows are affected.
func (r *AccountRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, accountID uuid.UUID, delta float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE account_id = $2 AND balance + $1 >= 0`
	result, err := tx.ExecContext(ctx, query, delta, accountID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Distinguish: does the account exist at all?
		var exists int
		err := tx.QueryRowContext(ctx, `SELECT 1 FROM accounts WHERE account_id = $1`, accountID).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrAccountNotFound
		}
		if err != nil {
			return err
		}
		return domain.ErrInsufficientBalance
	}

	return nil
}
