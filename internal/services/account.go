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

type AccountService struct {
	accountRepo     repo.AccountReadWriter
	tempAccountRepo repo.TempAccountReadWriter
	log             *zerolog.Logger
}

func NewAccountService(
	accountRepo repo.AccountReadWriter,
	tempAccountRepo repo.TempAccountReadWriter,
	log *zerolog.Logger) api.AccountApi {
	return &AccountService{
		accountRepo:     accountRepo,
		tempAccountRepo: tempAccountRepo,
		log:             log,
	}
}

func (service *AccountService) AddTempAccount(ctx context.Context, email string, lifetime time.Duration) (*domain.TempAccount, time.Time, error) {
	ok, err := service.tempAccountRepo.TempAccountExists(ctx, email)
	if err != nil {
		return nil, time.Now(), err
	}

	if ok {
		return nil, time.Now(), domain.ErrTempAccountExist
	}

	ok, err = service.accountRepo.AccountExistsByEmail(ctx, email)
	if err != nil {
		return nil, time.Now(), err
	}

	if ok {
		return nil, time.Now(), domain.ErrAccountExist
	}

	tempAccount := &domain.TempAccount{
		Code: utils.MustNewCode(),
	}

	err = service.tempAccountRepo.AddTempAccount(ctx, email, tempAccount, lifetime)
	if err == nil {
		return tempAccount, time.Now().Add(lifetime), nil
	}

	switch {
	case errors.Is(err, domain.ErrAccountExist):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, time.Now(), err
}

func (service *AccountService) RemoveTempAccount(ctx context.Context, email string) error {
	err := service.tempAccountRepo.RemoveTempAccount(ctx, email)

	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrNoTempAccount):
	default:
		service.log.Error().Err(err).Send()
	}

	return err
}

func (service *AccountService) SaveTempAccount(ctx context.Context, email, code string) (*domain.Account, error) {
	tempAccount, err := service.tempAccountRepo.GetTempAccount(ctx, email)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoTempAccount):
		default:
			service.log.Error().Err(err).Send()
		}

		return nil, err
	}

	if code != tempAccount.Code {
		return nil, domain.ErrWrongCode
	}

	account := &domain.Account{
		Id:        utils.MustNewId(),
		Email:     email,
		Perms:     domain.PermDefault,
		CreatedAt: time.Now(),
	}

	err = service.accountRepo.AddAccount(ctx, account)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccountExist):
		default:
			service.log.Error().Err(err).Send()
		}

		return nil, err
	}

	err = service.tempAccountRepo.RemoveTempAccount(ctx, email)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoTempAccount):
		default:
			service.log.Error().Err(err).Send()
		}

		return nil, err
	}

	return account, nil
}

func (service *AccountService) GetAccountById(ctx context.Context, id string) (*domain.Account, error) {
	account, err := service.accountRepo.GetAccountById(ctx, id)

	if err == nil {
		return account, nil
	}

	switch {
	case errors.Is(err, domain.ErrNoAccount):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, err
}

func (service *AccountService) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	account, err := service.accountRepo.GetAccountByEmail(ctx, email)

	if err == nil {
		return account, nil
	}

	switch {
	case errors.Is(err, domain.ErrNoAccount):
	default:
		service.log.Error().Err(err).Send()
	}

	return nil, err
}

func (service *AccountService) AccountExistsById(ctx context.Context, id string) (bool, error) {
	exists, err := service.accountRepo.AccountExistsById(ctx, id)

	if err != nil {
		service.log.Error().Err(err).Send()
	}

	return exists, err
}

func (service *AccountService) AccountExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := service.accountRepo.AccountExistsByEmail(ctx, email)

	if err != nil {
		service.log.Error().Err(err).Send()
	}

	return exists, err
}
