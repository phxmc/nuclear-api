package repo

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

type TelegramReader interface {
	GetChatState(ctx context.Context, chatId int64) (domain.ChatState, bool)
}

type TelegramWriter interface {
	SetChatState(ctx context.Context, chatId int64, state domain.ChatState, ttl time.Duration) error
}

type TelegramReadWriter interface {
	TelegramReader
	TelegramWriter
}
