package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/orewaee/nuclear-api/internal/app/domain"
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

func (repo *TokenRepo) TokenExists(ctx context.Context, prefix, token string) (bool, error) {
	key := fmt.Sprintf("%s:%s", prefix, token)
	result, err := repo.client.Exists(ctx, key).Result()

	if err == nil {
		return result == 1, nil
	}

	return false, err
}

func (repo *TokenRepo) AddToken(ctx context.Context, prefix, token string, lifetime time.Duration) error {
	exists, err := repo.TokenExists(ctx, prefix, token)
	if err != nil {
		return err
	}

	if exists {
		return domain.ErrTokenExist
	}

	key := fmt.Sprintf("%s:%s", prefix, token)
	status, err := repo.client.Set(ctx, key, true, lifetime).Result()
	if err != nil {
		return err
	}

	if status != "OK" {
		return errors.New("failed to add token")
	}

	return nil
}

func (repo *TokenRepo) RemoveToken(ctx context.Context, prefix, token string) error {
	exists, err := repo.TokenExists(ctx, prefix, token)
	if err != nil {
		return err
	}

	if !exists {
		return domain.ErrNoToken
	}

	key := fmt.Sprintf("%s:%s", prefix, token)
	result, err := repo.client.Del(ctx, key).Result()
	if err == nil {
		return err
	}

	if result != 1 {
		return errors.New("failed to remove token")
	}

	return nil
}
