package services

import (
	"context"
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
	if err != nil {
		service.log.Error().Err(err).Send()
		return nil, err
	}

	return pass, nil
}

func (service *PassService) GetPassByAccountId(ctx context.Context, accountId string) (*domain.Pass, error) {
	pass, err := service.passRepo.GetPassByAccountId(ctx, accountId)
	if err != nil {
		service.log.Error().Err(err).Send()
		return nil, err
	}

	if pass.From == nil && pass.To == nil {
		return pass, nil
	}

	now := time.Now()

	if pass.From != nil && now.Before(*pass.From) {
		return nil, domain.ErrInvalidPass
	}

	if pass.To != nil && now.After(*pass.To) {
		return nil, domain.ErrInvalidPass
	}

	return pass, nil
}

func (service *PassService) AddPass(ctx context.Context, accountId string, from *time.Time, to *time.Time) (*domain.Pass, error) {
	pass := &domain.Pass{
		Id:        utils.MustNewId(),
		AccountId: accountId,
		From:      from,
		To:        to,
		Active:    true,
		CreatedAt: time.Now(),
	}

	err := service.passRepo.AddPass(ctx, pass)
	if err != nil {
		service.log.Error().Err(err).Send()
		return nil, err
	}

	return pass, nil
}
