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
	ResetChatState(ctx context.Context, chatId int64) error
}

type TelegramReadWriter interface {
	TelegramReader
	TelegramWriter
}

type TempTelegramReader interface {
	GetTempTelegram(ctx context.Context, telegramId int64) (*domain.TempTelegram, error)

	// TempTelegramExists returns the bool value of the existence of a temporary telegram with the specified telegramId.
	TempTelegramExists(ctx context.Context, telegramId int64) (bool, error)
}

type TempTelegramWriter interface {
	// AddTempTelegram adds the specified temporary telegram.
	//
	// May return domain.ErrTempTelegramExist.
	AddTempTelegram(ctx context.Context, telegramId int64, tempTelegram *domain.TempTelegram, lifetime time.Duration) error

	// RemoveTempTelegram removes the temporary telegram with the specified telegramId.
	//
	// May return domain.ErrNoTempTelegram.
	RemoveTempTelegram(ctx context.Context, telegramId int64) error
}

type TempTelegramReadWriter interface {
	TempTelegramReader
	TempTelegramWriter
}
