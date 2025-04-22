package repo

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

type AccountReader interface {
	// GetAccountById returns the account by id.
	//
	// May return domain.ErrNoAccount.
	GetAccountById(ctx context.Context, id string) (*domain.Account, error)

	// GetAccountByEmail returns the account by email.
	//
	// This can return domain.ErrNoAccount.
	GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error)

	// GetAccountByTelegramId returns the account by telegramId.
	//
	// This can return domain.ErrNoAccount.
	GetAccountByTelegramId(ctx context.Context, telegramId int64) (*domain.Account, error)

	// AccountExistsById returns the bool value of the existence of an account with the specified id.
	AccountExistsById(ctx context.Context, id string) (bool, error)

	// AccountExistsByEmail returns the bool value of the existence of an account with the specified email.
	AccountExistsByEmail(ctx context.Context, email string) (bool, error)

	// AccountExistsByTelegramId returns the bool value of the existence of an account with the specified telegramId.
	AccountExistsByTelegramId(ctx context.Context, telegramId int64) (bool, error)
}

type AccountWriter interface {
	// AddAccount adds the specified account.
	//
	// May return domain.ErrAccountExist.
	AddAccount(ctx context.Context, account *domain.Account) error

	SetAccountTelegramId(ctx context.Context, accountId string, telegramId int64) error
}

type AccountReadWriter interface {
	AccountReader
	AccountWriter
}

type TempAccountReader interface {
	// GetTempAccount returns the temporary account with the specified email.
	//
	// May return domain.ErrNoTempAccount.
	GetTempAccount(ctx context.Context, email string) (*domain.TempAccount, error)

	// TempAccountExists returns the bool value of the existence of a temporary account with the specified email.
	TempAccountExists(ctx context.Context, email string) (bool, error)
}

type TempAccountWriter interface {
	// AddTempAccount adds the specified temporary account.
	//
	// May return domain.ErrTempAccountExist.
	AddTempAccount(ctx context.Context, email string, tempAccount *domain.TempAccount, lifetime time.Duration) error

	// RemoveTempAccount removes the temporary account with the specified email.
	//
	// May return domain.ErrNoTempAccount.
	RemoveTempAccount(ctx context.Context, email string) error
}

type TempAccountReadWriter interface {
	TempAccountReader
	TempAccountWriter
}
