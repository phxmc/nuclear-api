package services

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/rs/zerolog"
	"time"
)

type TelegramService struct {
	telegramRepo repo.TelegramReadWriter
	log          *zerolog.Logger
}

func NewTelegramService(telegramRepo repo.TelegramReadWriter, log *zerolog.Logger) api.TelegramApi {
	return &TelegramService{
		telegramRepo: telegramRepo,
		log:          log,
	}
}

func (service *TelegramService) GetChatState(ctx context.Context, chatId int64) (domain.ChatState, bool) {
	return service.telegramRepo.GetChatState(ctx, chatId)
}

func (service *TelegramService) SetChatState(ctx context.Context, chatId int64, state domain.ChatState, ttl time.Duration) error {
	err := service.telegramRepo.SetChatState(ctx, chatId, state, ttl)
	if err != nil {
		service.log.Error().Err(err).Send()
		return err
	}

	return nil
}
