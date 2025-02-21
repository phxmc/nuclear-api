package repo

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

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
