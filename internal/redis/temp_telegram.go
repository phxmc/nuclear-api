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

type TempTelegramRepo struct {
	client *goredis.Client
	prefix string
}

func NewTempTelegramRepo(client *goredis.Client) repo.TempTelegramReadWriter {
	return &TempTelegramRepo{client, "temp_telegram"}
}

func (repo *TempTelegramRepo) GetTempTelegram(ctx context.Context, telegramId int64) (*domain.TempTelegram, error) {
	key := fmt.Sprintf("%s:%d", repo.prefix, telegramId)
	data, err := repo.client.Get(ctx, key).Result()

	if err != nil {
		switch {
		case errors.Is(err, goredis.Nil):
			return nil, domain.ErrNoTempTelegram
		default:
			return nil, err
		}
	}

	var buffer bytes.Buffer
	buffer.WriteString(data)
	decoder := gob.NewDecoder(&buffer)

	tempTelegram := new(domain.TempTelegram)
	if err := decoder.Decode(tempTelegram); err != nil {
		return nil, err
	}

	return tempTelegram, nil
}

func (repo *TempTelegramRepo) TempTelegramExists(ctx context.Context, telegramId int64) (bool, error) {
	key := fmt.Sprintf("%s:%d", repo.prefix, telegramId)
	exists, err := repo.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (repo *TempTelegramRepo) AddTempTelegram(ctx context.Context, telegramId int64, tempTelegram *domain.TempTelegram, lifetime time.Duration) error {
	exists, err := repo.TempTelegramExists(ctx, telegramId)
	if err != nil {
		return err
	}

	if exists {
		return domain.ErrTempTelegramExist
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(tempTelegram); err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%d", repo.prefix, telegramId)
	return repo.client.Set(ctx, key, buffer.Bytes(), lifetime).Err()
}

func (repo *TempTelegramRepo) RemoveTempTelegram(ctx context.Context, telegramId int64) error {
	exists, err := repo.TempTelegramExists(ctx, telegramId)
	if err != nil {
		return err
	}

	if !exists {
		return domain.ErrNoTempTelegram
	}

	key := fmt.Sprintf("%s:%d", repo.prefix, telegramId)
	return repo.client.Del(ctx, key).Err()
}
