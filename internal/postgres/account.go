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
	account := new(domain.Account)
	err := withTx(ctx, repo.pool, func(tx pgx.Tx) error {
		err := tx.
			QueryRow(ctx, "SELECT * FROM accounts WHERE id = $1", id).
			Scan(&account.Id, &account.Email, &account.Perms, &account.CreatedAt)

		if err != nil && errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNoAccount
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo *AccountRepo) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	account := new(domain.Account)
	err := withTx(ctx, repo.pool, func(tx pgx.Tx) error {
		err := tx.
			QueryRow(ctx, "SELECT * FROM accounts WHERE email = $1", email).
			Scan(&account.Id, &account.Email, &account.Perms, &account.CreatedAt)

		if err != nil && errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNoAccount
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo *AccountRepo) AccountExistsById(ctx context.Context, id string) (bool, error) {
	exists := false
	err := withTx(ctx, repo.pool, func(tx pgx.Tx) error {
		return tx.
			QueryRow(ctx, "SElECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", id).
			Scan(&exists)
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *AccountRepo) AccountExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists := false
	err := withTx(ctx, repo.pool, func(tx pgx.Tx) error {
		return tx.
			QueryRow(ctx, "SElECT EXISTS(SELECT 1 FROM accounts WHERE email = $1)", email).
			Scan(&exists)
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *AccountRepo) AddAccount(ctx context.Context, account *domain.Account) error {
	err := withTx(ctx, repo.pool, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO accounts (id, email, perms, created_at) VALUES ($1, $2, $3, $4)",
			&account.Id, &account.Email, &account.Perms, &account.CreatedAt)

		return err
	})

	return err
}
