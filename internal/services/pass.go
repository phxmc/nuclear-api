package services

import (
	"context"
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/rs/zerolog"
	"time"
)

type PassService struct {
	passRepo repo.PassReadWriter
	log      *zerolog.Logger
}

func NewPassService(
	passRepo repo.PassReadWriter,
	log *zerolog.Logger) api.PassApi {
	return &PassService{
		passRepo: passRepo,
		log:      log,
	}
}

func (service *PassService) GetPassById(ctx context.Context, id string) (*domain.Pass, error) {
	pass, err := service.passRepo.GetPassById(ctx, id)

	if err == nil {
		return pass, nil
	}

	switch {
	case errors.Is(err, domain.ErrNoPass):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, err
}

func (service *PassService) GetPassByAccountId(ctx context.Context, accountId string) (*domain.Pass, error) {
	pass, err := service.passRepo.GetPassByAccountId(ctx, accountId)

	if err == nil {
		return pass, nil
	}

	switch {
	case errors.Is(err, domain.ErrNoAccount):
	case errors.Is(err, domain.ErrNoPass):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, err
}

func (service *PassService) GetPassHistoryByAccountId(ctx context.Context, accountId string) ([]*domain.Pass, error) {
	passes, err := service.passRepo.GetPassHistoryByAccountId(ctx, accountId)

	if err == nil {
		return passes, nil
	}

	switch {
	case errors.Is(err, domain.ErrNoAccount):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, err
}

func (service *PassService) SetPass(ctx context.Context, accountId string, from *time.Time, to *time.Time) (*domain.Pass, error) {
	pass := &domain.Pass{
		Id:        utils.MustNewId(),
		From:      from,
		To:        to,
		CreatedAt: time.Now(),
	}

	err := service.passRepo.SetPass(ctx, accountId, pass)

	if err == nil {
		return pass, nil
	}

	switch {
	case errors.Is(err, domain.ErrNoAccount):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, err
}
