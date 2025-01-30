package repo

import (
	"context"
	"time"
)

// TokenReader contains methods for reading tokens
type TokenReader interface {
	TokenExists(ctx context.Context, token string) (bool, error)
}

// TokenWriter contains methods for writing tokens
type TokenWriter interface {
	AddToken(ctx context.Context, token string, lifetime time.Duration) error
	RemoveToken(ctx context.Context, token string) error
}

// TokenReadWriter is a wrapper for TokenReader and TokenWriter
type TokenReadWriter interface {
	TokenReader
	TokenWriter
}
