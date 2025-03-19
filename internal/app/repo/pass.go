package repo

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
)

type PassReader interface {
	// GetPassById returns pass with the specified id.
	//
	// May return domain.ErrNoPass.
	GetPassById(ctx context.Context, id string) (*domain.Pass, error)

	// GetPassByAccountId returns the active pass for the specified account.
	//
	// May return domain.ErrNoAccount, domain.ErrNoPass, domain.ErrInvalidPass.
	GetPassByAccountId(ctx context.Context, accountId string) (*domain.Pass, error)

	// GetPassHistoryByAccountId returns the pass history for the specified account in chronological order.
	// All passes in history are inactive.
	//
	// May return domain.ErrNoAccount.
	GetPassHistoryByAccountId(ctx context.Context, accountId string) ([]*domain.Pass, error)
}

type PassWriter interface {
	// SetPass creates a new pass and sets it to the specified account.
	// If an active pass was associated with the account, it will be marked as inactive.
	//
	// May return domain.ErrNoAccount.
	SetPass(ctx context.Context, accountId string, pass *domain.Pass) error
}

type PassReadWriter interface {
	PassReader
	PassWriter
}
