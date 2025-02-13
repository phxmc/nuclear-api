package api

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

type PassApi interface {
	// GetPassById
	//
	// May return domain.ErrNoPass
	GetPassById(ctx context.Context, id string) (*domain.Pass, error)

	// GetPassByAccountId returns the active pass with a valid period for the specified account
	//
	// May return domain.ErrNoPass
	GetPassByAccountId(ctx context.Context, accountId string) (*domain.Pass, error)

	// AddPass
	//
	// Returns domain.ErrNoAccount if there is no account with the specified accountId
	//
	// Returns domain.ErrPassExist if an active pass is already bound to the account with the specified accountId
	AddPass(ctx context.Context, accountId string, from *time.Time, to *time.Time) (*domain.Pass, error)
}
