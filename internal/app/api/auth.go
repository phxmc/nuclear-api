package api

import (
	"context"
	"time"
)

type AuthApi interface {
	// Login returns the generated temporary code and its deadline.
	//
	// May return domain.ErrAccountNotExist, domain.ErrLoginCodeExist.
	Login(ctx context.Context, email string) (string, time.Time, error)

	// LoginCode returns a pair of tokens.
	//
	// May return domain.ErrNoLoginCode, domain.ErrWrongCode.
	LoginCode(ctx context.Context, email, code string) (string, string, error)

	// GenerateToken generates a signed token with the specified claims.
	GenerateToken(claims map[string]interface{}, key string) (string, error)

	// WhitelistToken whitelists the specified token with a prefix.
	WhitelistToken(ctx context.Context, prefix, token string, lifetime time.Duration) error

	// RevokeToken revokes the specified token with a prefix.
	RevokeToken(ctx context.Context, prefix, token string) error

	// GetTokenClaims retrieves claims from the token.
	//
	// May return domain.ErrInvalidToken.
	GetTokenClaims(token string, key string) (map[string]interface{}, error)

	// RefreshTokens returns a pair of tokens.
	//
	// May return domain.ErrInvalidToken.
	RefreshTokens(ctx context.Context, prefix, token string) (string, string, error)
}
