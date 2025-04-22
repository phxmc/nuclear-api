package api

import (
	"context"
	"time"

	"github.com/orewaee/nuclear-api/internal/app/domain"
)

type TelegramApi interface {
	GetChatState(ctx context.Context, chatId int64) (domain.ChatState, bool)
	GetTempTelegram(ctx context.Context, telegramId int64) (*domain.TempTelegram, error)
	TempTelegramExists(ctx context.Context, telegramId int64) (bool, error)

	SetChatState(ctx context.Context, chatId int64, state domain.ChatState, ttl time.Duration) error
	ResetChatState(ctx context.Context, chatId int64) error
	AddTempTelegram(ctx context.Context, telegramId int64, code, email string, lifetime time.Duration) error
	RemoveTempTelegram(ctx context.Context, telegramId int64) error

	SendMessage(ctx context.Context, chatId int64, markdown string) error
}
