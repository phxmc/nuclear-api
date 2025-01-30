package api

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

type AccountApi interface {
	// AddTempAccount
	//
	// This can return domain.ErrTempAccountExist, domain.ErrAccountExist
	AddTempAccount(ctx context.Context, email string, lifetime time.Duration) (*domain.TempAccount, time.Time, error)

	// RemoveTempAccount
	//
	// This can only return internal errors
	RemoveTempAccount(ctx context.Context, email string) error

	// SaveTempAccount
	//
	// This can return domain.ErrWrongCode, domain.ErrTempAccountNotExist
	SaveTempAccount(ctx context.Context, email, code string) (*domain.Account, error)

	// GetAccountById
	//
	// This can return domain.ErrAccountNotExist
	GetAccountById(ctx context.Context, id string) (*domain.Account, error)

	// GetAccountByEmail
	//
	// This can return domain.ErrAccountNotExist
	GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error)
}
