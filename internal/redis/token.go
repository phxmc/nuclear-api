package redis

import (
	"context"
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	goredis "github.com/redis/go-redis/v9"
	"time"
)

type TokenRepo struct {
	client *goredis.Client
}

func NewTokenRepo(client *goredis.Client) repo.TokenReadWriter {
	return &TokenRepo{client}
}

func (repo *TokenRepo) TokenExists(ctx context.Context, token string) (bool, error) {
	_, err := repo.client.Get(ctx, token).Result()

	if err != nil && errors.Is(err, goredis.Nil) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *TokenRepo) AddToken(ctx context.Context, token string, lifetime time.Duration) error {
	status, err := repo.client.Set(ctx, token, true, lifetime).Result()

	if err != nil {
		return err
	}

	if status != "OK" {
		return errors.New("failed to add token")
	}

	return nil
}

func (repo *TokenRepo) RemoveToken(ctx context.Context, token string) error {
	_, err := repo.client.Del(ctx, token).Result()

	if err != nil {
		return err
	}

	return nil
}
