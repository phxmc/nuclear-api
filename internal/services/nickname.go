package services

import (
	"context"
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/typedenv"
	"github.com/rs/zerolog"
	"time"
)

type NicknameService struct {
	nicknameRepo repo.NicknameReadWriter
	log          *zerolog.Logger
}

func NewNicknameService(
	nicknameRepo repo.NicknameReadWriter,
	log *zerolog.Logger) api.NicknameApi {
	return &NicknameService{
		nicknameRepo: nicknameRepo,
		log:          log,
	}
}

func (service *NicknameService) GetNicknameByAccountId(ctx context.Context, accountId string) (*domain.Nickname, error) {
	nickname, err := service.nicknameRepo.GetNicknameByAccountId(ctx, accountId)

	if err == nil {
		return nickname, nil
	}

	switch {
	case errors.Is(err, domain.ErrNoAccount):
	case errors.Is(err, domain.ErrNoNickname):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, err
}

func (service *NicknameService) GetNicknameHistoryByAccountId(ctx context.Context, accountId string) ([]*domain.Nickname, error) {
	nicknames, err := service.nicknameRepo.GetNicknameHistoryByAccountId(ctx, accountId)

	if err == nil {
		return nicknames, nil
	}

	switch {
	case errors.Is(err, domain.ErrNoAccount):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, err
}

func (service *NicknameService) NicknameExists(ctx context.Context, nickname string) (bool, error) {
	exists, err := service.nicknameRepo.NicknameExists(ctx, nickname)
	if err != nil {
		service.log.Error().Err(err).Send()
		return false, nil
	}

	return exists, nil
}

func (service *NicknameService) SetNickname(ctx context.Context, accountId, nickname string) (*domain.Nickname, error) {
	hasNickname := true
	oldNickname, err := service.GetNicknameByAccountId(ctx, accountId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoAccount):
			return nil, err
		case errors.Is(err, domain.ErrNoNickname):
			hasNickname = false
			err = nil
		default:
			service.log.Error().Err(err).Send()
			return nil, err
		}
	}

	if hasNickname {
		now := time.Now()
		cooldown := typedenv.Duration("NICKNAME_COOLDOWN", time.Hour*24*7)
		deadline := oldNickname.CreatedAt.Add(cooldown)
		if deadline.After(now) {
			return nil, domain.ErrNicknameCooldown
		}
	}

	exists, err := service.nicknameRepo.NicknameExists(ctx, nickname)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrNicknameExist
	}

	newNickname := &domain.Nickname{
		Value:     nickname,
		CreatedAt: time.Now(),
	}

	err = service.nicknameRepo.SetNickname(ctx, accountId, newNickname)
	if err == nil {
		return newNickname, nil
	}

	switch {
	case errors.Is(err, domain.ErrNoAccount):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, err
}
