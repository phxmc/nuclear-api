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
	row := repo.pool.QueryRow(ctx, "SELECT * FROM accounts WHERE id = $1", id)

	account := new(domain.Account)
	err := row.Scan(&account.Id, &account.Email, &account.Perms)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNoAccount
	}

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo *AccountRepo) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	row := repo.pool.QueryRow(ctx, "SELECT * FROM accounts WHERE email = $1", email)

	account := new(domain.Account)
	err := row.Scan(&account.Id, &account.Email, &account.Perms)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNoAccount
	}

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo *AccountRepo) AccountExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	row := repo.pool.QueryRow(ctx, "SELECT COUNT(*) FROM accounts WHERE email = $1", email)
	err := row.Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *AccountRepo) AddAccount(ctx context.Context, account *domain.Account) error {
	sql := "INSERT INTO accounts (id, email, perms) VALUES ($1, $2, $3)"
	_, err := repo.pool.Exec(ctx, sql, account.Id, account.Email, account.Perms)

	return err
}
