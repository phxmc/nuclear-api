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
	prefix string
}

func NewLoginCodeRepo(client *goredis.Client) repo.LoginCodeReadWriter {
	return &LoginCodeRepo{client, "login_code"}
}

func (repo *LoginCodeRepo) GetLoginCode(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("%s:%s", repo.prefix, email)
	data, err := repo.client.Get(ctx, key).Result()

	if err != nil {
		switch {
		case errors.Is(err, goredis.Nil):
			return "", domain.ErrLoginCodeNotExist
		default:
			return "", err
		}
	}

	return data, nil
}

func (repo *LoginCodeRepo) LoginCodeExists(ctx context.Context, email string) (bool, error) {
	key := fmt.Sprintf("%s:%s", repo.prefix, email)
	exists, err := repo.client.Exists(ctx, key).Result()

	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (repo *LoginCodeRepo) AddLoginCode(ctx context.Context, email string, code string, lifetime time.Duration) error {
	key := fmt.Sprintf("%s:%s", repo.prefix, email)
	return repo.client.Set(ctx, key, code, lifetime).Err()
}

func (repo *LoginCodeRepo) RemoveLoginCode(ctx context.Context, email string) error {
	key := fmt.Sprintf("%s:%s", repo.prefix, email)
	return repo.client.Del(ctx, key).Err()
}
