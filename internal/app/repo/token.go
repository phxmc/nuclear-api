package repo

import (
	"context"
	"time"
)

type TokenReader interface {
	// TokenExists checks for the existence of the specified token with a prefix.
	TokenExists(ctx context.Context, prefix, token string) (bool, error)
}

type TokenWriter interface {
	// AddToken adds the specified token with a prefix.
	//
	// May return domain.ErrTokenExist.
	AddToken(ctx context.Context, prefix, token string, lifetime time.Duration) error

	// RemoveToken removes the specified token with a prefix.
	//
	// May return domain.ErrNoToken.
	RemoveToken(ctx context.Context, prefix, token string) error
}

type TokenReadWriter interface {
	TokenReader
	TokenWriter
}
