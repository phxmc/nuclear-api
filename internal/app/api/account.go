package api

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

type AccountApi interface {
	// AddTempAccount creates a temporary account and maps it to the specified email.
	//
	// Returns the created temporary account and its deadline.
	//
	// May return domain.ErrTempAccountExist, domain.ErrAccountExist.
	AddTempAccount(ctx context.Context, email string, lifetime time.Duration) (*domain.TempAccount, time.Time, error)

	// RemoveTempAccount removes the temporary account with the specified email.
	RemoveTempAccount(ctx context.Context, email string) error

	// SaveTempAccount turns a temporary account into a permanent account.
	//
	// May return domain.ErrWrongCode, domain.ErrNoTempAccount.
	SaveTempAccount(ctx context.Context, email, code string) (*domain.Account, error)

	// GetAccountById returns the account with the specified id.
	//
	// May return domain.ErrNoAccount.
	GetAccountById(ctx context.Context, id string) (*domain.Account, error)

	// GetAccountByEmail returns the account with the specified email.
	//
	// May return domain.ErrNoAccount.
	GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error)

	// AccountExistsByEmail returns the bool value of the existence of an account with the specified email.
	AccountExistsByEmail(ctx context.Context, email string) (bool, error)
}
