package repo

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
)

type AccountReader interface {
	GetAccountById(ctx context.Context, id string) (*domain.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error)

	// AccountExistsByEmail returns true if an account with the specified email exists, otherwise false
	AccountExistsByEmail(ctx context.Context, email string) (bool, error)
}

type AccountWriter interface {
	AddAccount(ctx context.Context, account *domain.Account) error
}

type AccountReadWriter interface {
	AccountReader
	AccountWriter
}
