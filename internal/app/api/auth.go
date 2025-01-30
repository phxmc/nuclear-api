package api

import (
	"context"
	"time"
)

type AuthApi interface {
	// Login
	//
	// This can return domain.ErrAccountNotExist, domain.ErrLoginCodeExist
	Login(ctx context.Context, email string) (string, time.Time, error)

	// LoginCode
	//
	// This can return domain.ErrLoginCodeNotExist, domain.ErrWrongCode
	LoginCode(ctx context.Context, email, code string) (string, string, error)

	CreateToken(claims map[string]interface{}, key string) (string, error)

	WhitelistToken(ctx context.Context, refreshToken string, lifetime time.Duration) error

	GetTokenClaims(token string, key string) (map[string]interface{}, error)

	// RefreshToken
	//
	// This can return domain.ErrInvalidToken
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}
