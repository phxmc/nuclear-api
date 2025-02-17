package api

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

type PassApi interface {
	// GetPassById returns pass with the specified id.
	//
	// May return domain.ErrNoPass.
	GetPassById(ctx context.Context, id string) (*domain.Pass, error)

	// GetPassByAccountId returns the active pass for the specified account.
	//
	// May return domain.ErrNoAccount, domain.ErrNoPass.
	GetPassByAccountId(ctx context.Context, accountId string) (*domain.Pass, error)

	// GetPassHistoryByAccountId returns the pass history for the specified account in chronological order.
	// All passes in history are inactive.
	//
	// May return domain.ErrNoAccount.
	GetPassHistoryByAccountId(ctx context.Context, accountId string) ([]*domain.Pass, error)

	// SetPass creates a new pass and sets it to the specified account.
	// If an active pass was associated with the account, it will be marked as inactive.
	//
	// May return domain.ErrNoAccount.
	SetPass(ctx context.Context, accountId string, from *time.Time, to *time.Time) (*domain.Pass, error)
}
