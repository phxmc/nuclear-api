package api

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

type TelegramApi interface {
	GetChatState(ctx context.Context, chatId int64) (domain.ChatState, bool)
	SetChatState(ctx context.Context, chatId int64, state domain.ChatState, ttl time.Duration) error
}
