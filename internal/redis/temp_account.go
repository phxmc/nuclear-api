package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	goredis "github.com/redis/go-redis/v9"
	"time"
)

type TempAccountRepo struct {
	client *goredis.Client
	prefix string
}

func NewTempAccountRepo(client *goredis.Client) repo.TempAccountReadWriter {
	return &TempAccountRepo{client, "temp_account"}
}

func (repo *TempAccountRepo) GetTempAccount(ctx context.Context, email string) (*domain.TempAccount, error) {
	key := fmt.Sprintf("%s:%s", repo.prefix, email)
	data, err := repo.client.Get(ctx, key).Result()

	if err != nil {
		switch {
		case errors.Is(err, goredis.Nil):
			return nil, domain.ErrTempAccountNotExist
		default:
			return nil, err
		}
	}

	var buffer bytes.Buffer
	buffer.WriteString(data)
	decoder := gob.NewDecoder(&buffer)

	tempAccount := new(domain.TempAccount)
	if err := decoder.Decode(tempAccount); err != nil {
		return nil, err
	}

	return tempAccount, nil
}

func (repo *TempAccountRepo) TempAccountExists(ctx context.Context, email string) (bool, error) {
	key := fmt.Sprintf("%s:%s", repo.prefix, email)
	exists, err := repo.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (repo *TempAccountRepo) AddTempAccount(ctx context.Context, email string, tempAccount *domain.TempAccount, lifetime time.Duration) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(tempAccount); err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", repo.prefix, email)
	return repo.client.Set(ctx, key, buffer.Bytes(), lifetime).Err()
}

func (repo *TempAccountRepo) RemoveTempAccount(ctx context.Context, email string) error {
	key := fmt.Sprintf("%s:%s", repo.prefix, email)
	return repo.client.Del(ctx, key).Err()
}
