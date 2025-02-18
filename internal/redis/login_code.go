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

type LoginCodeRepo struct {
	client *goredis.Client
}

func NewLoginCodeRepo(client *goredis.Client) repo.LoginCodeReadWriter {
	return &LoginCodeRepo{client}
}

func (repo *LoginCodeRepo) GetLoginCode(ctx context.Context, prefix, email string) (string, error) {
	key := fmt.Sprintf("%s:%s", prefix, email)
	code, err := repo.client.Get(ctx, key).Result()

	if err == nil {
		return code, nil
	}

	switch {
	case errors.Is(err, goredis.Nil):
		return "", domain.ErrNoLoginCode
	default:
		return "", err
	}
}

func (repo *LoginCodeRepo) LoginCodeExists(ctx context.Context, prefix, email string) (bool, error) {
	key := fmt.Sprintf("%s:%s", prefix, email)
	result, err := repo.client.Exists(ctx, key).Result()

	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (repo *LoginCodeRepo) AddLoginCode(ctx context.Context, prefix, email, code string, lifetime time.Duration) error {
	exists, err := repo.LoginCodeExists(ctx, prefix, email)
	if err != nil {
		return err
	}

	if exists {
		return domain.ErrLoginCodeExist
	}

	key := fmt.Sprintf("%s:%s", prefix, email)
	return repo.client.Set(ctx, key, code, lifetime).Err()
}

func (repo *LoginCodeRepo) RemoveLoginCode(ctx context.Context, prefix, email string) error {
	exists, err := repo.LoginCodeExists(ctx, prefix, email)
	if err != nil {
		return err
	}

	if !exists {
		return domain.ErrNoLoginCode
	}

	key := fmt.Sprintf("%s:%s", prefix, email)
	return repo.client.Del(ctx, key).Err()
}
