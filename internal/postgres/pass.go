package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
)

type PassRepo struct {
	pool *pgxpool.Pool
}

func NewPassRepo(pool *pgxpool.Pool) repo.PassReadWriter {
	return &PassRepo{pool: pool}
}

func (repo *PassRepo) GetPassById(ctx context.Context, id string) (*domain.Pass, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	pass := new(domain.Pass)

	row := tx.QueryRow(ctx, "SELECT * FROM passes WHERE id = $1", id)
	err = row.Scan(&pass.Id, &pass.AccountId, &pass.From, &pass.To, &pass.Active, &pass.CreatedAt)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNoPass
	}

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return pass, nil
}

func (repo *PassRepo) GetPassByAccountId(ctx context.Context, accountId string) (*domain.Pass, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	pass := new(domain.Pass)

	row := tx.QueryRow(ctx, "SELECT * FROM passes WHERE account_id = $1 AND active = true", accountId)
	err = row.Scan(&pass.Id, &pass.AccountId, &pass.From, &pass.To, &pass.Active, &pass.CreatedAt)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNoPass
	}

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return pass, nil
}

func (repo *PassRepo) AddPass(ctx context.Context, pass *domain.Pass) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, "INSERT INTO passes (id, account_id, from, to, active, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		pass.Id, pass.AccountId, pass.From, pass.To, pass.Active, pass.CreatedAt)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
