package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
)

type NicknameRepo struct {
	pool *pgxpool.Pool
}

func NewNicknameRepo(pool *pgxpool.Pool) repo.NicknameReadWriter {
	return &NicknameRepo{pool}
}

func (repo *NicknameRepo) GetNicknameByAccountId(ctx context.Context, accountId string) (*domain.Nickname, error) {
	nickname := new(domain.Nickname)
	err := withTx(ctx, repo.pool, func(tx pgx.Tx) error {
		exists := false
		err := tx.
			QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", accountId).
			Scan(&exists)

		if err != nil {
			return err
		}

		if !exists {
			return domain.ErrNoAccount
		}

		sql := `
			SELECT nickname, created_at
			FROM account_nicknames
			WHERE account_id = $1 AND is_active = TRUE
		`

		err = tx.
			QueryRow(ctx, sql, accountId).
			Scan(&nickname.Value, &nickname.CreatedAt)

		if err != nil && errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNoNickname
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return nickname, nil
}

func (repo *NicknameRepo) GetNicknameHistoryByAccountId(ctx context.Context, accountId string) ([]*domain.Nickname, error) {
	nicknames := make([]*domain.Nickname, 0)
	err := withTx(ctx, repo.pool, func(tx pgx.Tx) error {
		exists := false
		err := tx.
			QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", accountId).
			Scan(&exists)

		if err != nil {
			return err
		}

		if !exists {
			return domain.ErrNoAccount
		}

		sql := `
			SELECT nickname, created_at
			FROM account_nicknames
			WHERE account_id = $1 AND is_active = FALSE
			ORDER BY created_at DESC
		`

		rows, err := tx.Query(ctx, sql, accountId)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		nicknames, err = pgx.CollectRows[*domain.Nickname](rows, func(row pgx.CollectableRow) (*domain.Nickname, error) {
			nickname := new(domain.Nickname)

			err := row.Scan(&nickname.Value, &nickname.CreatedAt)
			if err != nil {
				return nil, err
			}

			return nickname, nil
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return nicknames, nil
}

func (repo *NicknameRepo) NicknameExists(ctx context.Context, nickname string) (bool, error) {
	exists := false
	err := withTx(ctx, repo.pool, func(tx pgx.Tx) error {
		err := tx.
			QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM account_nicknames WHERE nickname = $1 AND is_active = TRUE)", nickname).
			Scan(&exists)

		return err
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *NicknameRepo) SetNickname(ctx context.Context, accountId string, nickname *domain.Nickname) error {
	return withTx(ctx, repo.pool, func(tx pgx.Tx) error {
		exists := false
		err := tx.
			QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", accountId).
			Scan(&exists)

		if err != nil {
			return err
		}

		if !exists {
			return domain.ErrNoAccount
		}

		_, err = tx.Exec(
			ctx, `UPDATE account_nicknames SET is_active = FALSE WHERE account_id = $1 AND is_active = TRUE`,
			accountId)

		if err != nil {
			return err
		}

		_, err = tx.Exec(
			ctx, `INSERT INTO account_nicknames (account_id, nickname, is_active, created_at) VALUES ($1, $2, TRUE, $3)`,
			accountId, nickname.Value, nickname.CreatedAt)

		return err
	})
}
