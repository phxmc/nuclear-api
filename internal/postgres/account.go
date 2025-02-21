package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
)

type AccountRepo struct {
	pool *pgxpool.Pool
}

func NewAccountRepo(pool *pgxpool.Pool) repo.AccountReadWriter {
	return &AccountRepo{pool}
}

func (repo *AccountRepo) GetAccountById(ctx context.Context, id string) (*domain.Account, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	account := new(domain.Account)
	err = tx.
		QueryRow(ctx, "SELECT * FROM accounts WHERE id = $1", id).
		Scan(&account.Id, &account.Email, &account.Perms, &account.CreatedAt)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNoAccount
	}

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo *AccountRepo) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	account := new(domain.Account)
	err = tx.
		QueryRow(ctx, "SELECT * FROM accounts WHERE email = $1", email).
		Scan(&account.Id, &account.Email, &account.Perms, &account.CreatedAt)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNoAccount
	}

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo *AccountRepo) AccountExistsByEmail(ctx context.Context, email string) (bool, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return false, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	exists := false
	err = tx.
		QueryRow(ctx, "SElECT EXISTS(SELECT 1 FROM accounts WHERE email = $1)", email).
		Scan(&exists)

	if err != nil {
		return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *AccountRepo) AddAccount(ctx context.Context, account *domain.Account) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, "INSERT INTO accounts (id, email, perms, created_at) VALUES ($1, $2, $3, $4)",
		&account.Id, &account.Email, &account.Perms, &account.CreatedAt)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
