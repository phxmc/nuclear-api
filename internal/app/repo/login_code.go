package repo

import (
	"context"
	"time"
)

type LoginCodeReader interface {
	// GetLoginCode may return domain.ErrLoginCodeNotExist
	GetLoginCode(ctx context.Context, email string) (string, error)

	LoginCodeExists(ctx context.Context, email string) (bool, error)
}

type LoginCodeWriter interface {
	AddLoginCode(ctx context.Context, email string, code string, lifetime time.Duration) error

	RemoveLoginCode(ctx context.Context, email string) error
}

type LoginCodeReadWriter interface {
	LoginCodeReader
	LoginCodeWriter
}
