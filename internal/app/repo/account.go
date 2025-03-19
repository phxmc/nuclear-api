package repo

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
)

type AccountReader interface {
	// GetAccountById returns the account by id
	//
	// May return domain.ErrNoAccount
	GetAccountById(ctx context.Context, id string) (*domain.Account, error)

	// GetAccountByEmail returns the account by email
	//
	// This can return domain.ErrNoAccount
	GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error)

	// AccountExistsById returns the bool value of the existence of an account with the specified id.
	AccountExistsById(ctx context.Context, id string) (bool, error)

	// AccountExistsByEmail returns the bool value of the existence of an account with the specified email.
	AccountExistsByEmail(ctx context.Context, email string) (bool, error)
}

type AccountWriter interface {
	// AddAccount adds the specified account.
	//
	// May return domain.ErrAccountExist.
	AddAccount(ctx context.Context, account *domain.Account) error
}

type AccountReadWriter interface {
	AccountReader
	AccountWriter
}
