package postgres

import (
	"context"
	"database/sql"
	"errors"

	"AccountFlow/internal/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *accountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `INSERT INTO accounts (account_id, document_number, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, account.AccountID, account.DocumentNumber, account.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrDocumentAlreadyUsed
		}
		return err
	}
	return nil
}

func (r *accountRepository) FindByID(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	query := `SELECT account_id, document_number, created_at FROM accounts WHERE account_id = $1`
	var a domain.Account
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(&a.AccountID, &a.DocumentNumber, &a.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}
