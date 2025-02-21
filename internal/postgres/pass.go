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
	return &PassRepo{pool}
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

	err = tx.
		QueryRow(ctx, "SELECT * FROM passes WHERE id = $1", id).
		Scan(&pass.Id, &pass.From, &pass.To, &pass.CreatedAt)

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

	exists := false

	err = tx.
		QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", accountId).
		Scan(&exists)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, domain.ErrNoAccount
	}

	pass := new(domain.Pass)

	sql := `
		SELECT p.id, p."from", p."to", p.created_at
		FROM account_passes ap
		JOIN passes p ON ap.pass_id = p.id
		WHERE ap.account_id = $1 AND ap.is_active = TRUE
	`

	err = tx.
		QueryRow(ctx, sql, accountId).
		Scan(&pass.Id, &pass.From, &pass.To, &pass.CreatedAt)

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

func (repo *PassRepo) GetPassHistoryByAccountId(ctx context.Context, accountId string) ([]*domain.Pass, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	exists := false

	err = tx.
		QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", accountId).
		Scan(&exists)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, domain.ErrNoAccount
	}

	passes := make([]*domain.Pass, 0)

	sql := `
		SELECT p.id, p."from", p."to", p.created_at
		FROM account_passes ap
		JOIN passes p ON ap.pass_id = p.id
		WHERE ap.account_id = $1 AND ap.is_active = FALSE
		ORDER BY p.created_at DESC
	`

	rows, err := tx.Query(ctx, sql, accountId)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return passes, nil
	}

	if err != nil {
		return nil, err
	}

	passes, err = pgx.CollectRows[*domain.Pass](rows, func(row pgx.CollectableRow) (*domain.Pass, error) {
		pass := new(domain.Pass)

		err := row.Scan(&pass.Id, &pass.From, &pass.To, &pass.CreatedAt)
		if err != nil {
			return nil, err
		}

		return pass, nil
	})

	if err != nil {
		return nil, err
	}

	return passes, err
}

func (repo *PassRepo) SetPass(ctx context.Context, accountId string, pass *domain.Pass) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	exists := false

	err = tx.
		QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", accountId).
		Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return domain.ErrNoAccount
	}

	_, err = tx.Exec(
		ctx, `INSERT INTO passes (id, "from", "to", created_at) VALUES ($1, $2, $3, $4)`,
		pass.Id, pass.From, pass.To, pass.CreatedAt)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx, `UPDATE account_passes SET is_active = FALSE WHERE account_id = $1 AND is_active = TRUE`,
		accountId)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx, `INSERT INTO account_passes (account_id, pass_id, is_active) VALUES ($1, $2, TRUE)`,
		accountId, pass.Id)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
