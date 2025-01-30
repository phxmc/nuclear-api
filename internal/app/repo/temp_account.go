package repo

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

type TempAccountReader interface {
	// GetTempAccount returns the *domain.Account associated with the specified email address
	//
	// This can return domain.ErrTempAccountNotExist
	GetTempAccount(ctx context.Context, email string) (*domain.TempAccount, error)

	// TempAccountExists checks by email to see if a temp account exists
	TempAccountExists(ctx context.Context, email string) (bool, error)
}

type TempAccountWriter interface {
	AddTempAccount(ctx context.Context, email string, tempAccount *domain.TempAccount, lifetime time.Duration) error

	RemoveTempAccount(ctx context.Context, email string) error
}

type TempAccountReadWriter interface {
	TempAccountReader
	TempAccountWriter
}
