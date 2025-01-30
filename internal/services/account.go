package services

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/utils"
	"time"
)

type AccountService struct {
	accountRepo     repo.AccountReadWriter
	tempAccountRepo repo.TempAccountReadWriter
}

func NewAccountService(
	accountRepo repo.AccountReadWriter,
	tempAccountRepo repo.TempAccountReadWriter) api.AccountApi {
	return &AccountService{
		accountRepo:     accountRepo,
		tempAccountRepo: tempAccountRepo,
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
		Code: utils.NewCode(),
	}

	err = service.tempAccountRepo.AddTempAccount(ctx, email, tempAccount, lifetime)
	if err != nil {
		return nil, time.Now(), err
	}

	return tempAccount, time.Now().Add(lifetime), nil
}

func (service *AccountService) RemoveTempAccount(ctx context.Context, email string) error {
	return service.tempAccountRepo.RemoveTempAccount(ctx, email)
}

func (service *AccountService) SaveTempAccount(ctx context.Context, email, code string) (*domain.Account, error) {
	tempAccount, err := service.tempAccountRepo.GetTempAccount(ctx, email)
	if err != nil {
		return nil, err
	}

	if code != tempAccount.Code {
		return nil, domain.ErrWrongCode
	}

	account := &domain.Account{
		Id:    utils.MustNewId(),
		Email: email,
		Perms: domain.PermDefault,
	}

	err = service.accountRepo.AddAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	err = service.tempAccountRepo.RemoveTempAccount(ctx, email)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (service *AccountService) GetAccountById(ctx context.Context, id string) (*domain.Account, error) {
	return service.accountRepo.GetAccountById(ctx, id)
}

func (service *AccountService) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	return service.accountRepo.GetAccountByEmail(ctx, email)
}
