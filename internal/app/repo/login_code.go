package repo

import (
	"context"
	"time"
)

type LoginCodeReader interface {
	// GetLoginCode returns the code at the specified email.
	//
	// May return domain.ErrNoLoginCode.
	GetLoginCode(ctx context.Context, prefix, email string) (string, error)

	// LoginCodeExists checks for the existence of the specified code.
	LoginCodeExists(ctx context.Context, prefix, email string) (bool, error)
}

type LoginCodeWriter interface {
	// AddLoginCode adds the specified code.
	//
	// May return domain.ErrLoginCodeExist.
	AddLoginCode(ctx context.Context, prefix, email, code string, lifetime time.Duration) error

	// RemoveLoginCode removes the specified code by email.
	//
	// May return domain.ErrNoLoginCode.
	RemoveLoginCode(ctx context.Context, prefix, email string) error
}

type LoginCodeReadWriter interface {
	LoginCodeReader
	LoginCodeWriter
}
