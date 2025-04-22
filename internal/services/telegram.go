package services

import (
	"context"
	"time"

	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/broker"
	"github.com/rs/zerolog"
)

type TelegramService struct {
	telegramRepo     repo.TelegramReadWriter
	tempTelegramRepo repo.TempTelegramReadWriter
	broker           *broker.Broker[*domain.Message]
	log              *zerolog.Logger
}

func NewTelegramService(
	telegramRepo repo.TelegramReadWriter,
	tempTelegramRepo repo.TempTelegramReadWriter,
	messageBroker *broker.Broker[*domain.Message],
	log *zerolog.Logger) api.TelegramApi {
	return &TelegramService{
		telegramRepo:     telegramRepo,
		tempTelegramRepo: tempTelegramRepo,
		broker:           messageBroker,
		log:              log,
	}
}

func (service *TelegramService) GetChatState(ctx context.Context, chatId int64) (domain.ChatState, bool) {
	return service.telegramRepo.GetChatState(ctx, chatId)
}

func (service *TelegramService) GetTempTelegram(ctx context.Context, telegramId int64) (*domain.TempTelegram, error) {
	tempTelegram, err := service.tempTelegramRepo.GetTempTelegram(ctx, telegramId)
	if err != nil {
		service.log.Error().Err(err).Send()
		return nil, err
	}

	return tempTelegram, nil
}

func (service *TelegramService) TempTelegramExists(ctx context.Context, telegramId int64) (bool, error) {
	exists, err := service.tempTelegramRepo.TempTelegramExists(ctx, telegramId)
	if err != nil {
		service.log.Error().Err(err).Send()
		return false, err
	}

	return exists, nil
}

func (service *TelegramService) SetChatState(ctx context.Context, chatId int64, state domain.ChatState, ttl time.Duration) error {
	err := service.telegramRepo.SetChatState(ctx, chatId, state, ttl)
	if err != nil {
		service.log.Error().Err(err).Send()
		return err
	}

	return nil
}

func (service *TelegramService) ResetChatState(ctx context.Context, chatId int64) error {
	err := service.telegramRepo.ResetChatState(ctx, chatId)
	if err != nil {
		service.log.Error().Err(err).Send()
		return err
	}

	return nil
}

func (service *TelegramService) AddTempTelegram(ctx context.Context, telegramId int64, code, email string, lifetime time.Duration) error {
	tempTelegram := &domain.TempTelegram{Code: code, Email: email}
	err := service.tempTelegramRepo.AddTempTelegram(ctx, telegramId, tempTelegram, lifetime)
	if err != nil {
		service.log.Error().Err(err).Send()
	}

	return err
}

func (service *TelegramService) RemoveTempTelegram(ctx context.Context, telegramId int64) error {
	err := service.tempTelegramRepo.RemoveTempTelegram(ctx, telegramId)
	if err != nil {
		service.log.Error().Err(err).Send()
	}

	return err
}

func (service *TelegramService) SendMessage(ctx context.Context, chatId int64, markdown string) error {
	message := &domain.Message{ChatId: chatId, Markdown: markdown}
	service.broker.Publish(message)
	return nil
}
