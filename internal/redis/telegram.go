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

type TelegramRepo struct {
	client *goredis.Client
}

func NewTelegramRepo(client *goredis.Client) repo.TelegramReadWriter {
	return &TelegramRepo{client}
}

func (repo *TelegramRepo) GetChatState(ctx context.Context, chatId int64) (domain.ChatState, bool) {
	key := fmt.Sprintf("chat_state:%d", chatId)
	val, err := repo.client.Get(ctx, key).Result()
	if errors.Is(err, goredis.Nil) {
		return "", false
	}

	chatState := domain.ChatState(val)
	if !chatState.Valid() {
		return "", false
	}

	return chatState, true
}

func (repo *TelegramRepo) SetChatState(ctx context.Context, chatId int64, state domain.ChatState, ttl time.Duration) error {
	key := fmt.Sprintf("chat_state:%d", chatId)
	return repo.client.Set(ctx, key, state, ttl).Err()
}
