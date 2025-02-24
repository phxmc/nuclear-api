package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/orewaee/typedenv"
	"github.com/rs/zerolog"
	"time"
)

type AuthService struct {
	accountRepo   repo.AccountReadWriter
	loginCodeRepo repo.LoginCodeReadWriter
	tokenRepo     repo.TokenReadWriter
	log           *zerolog.Logger
}

func NewAuthService(
	accountRepo repo.AccountReadWriter,
	loginCodeRepo repo.LoginCodeReadWriter,
	tokenRepo repo.TokenReadWriter,
	log *zerolog.Logger) api.AuthApi {
	return &AuthService{
		accountRepo:   accountRepo,
		loginCodeRepo: loginCodeRepo,
		tokenRepo:     tokenRepo,
		log:           log,
	}
}

func (service *AuthService) Login(ctx context.Context, email string) (string, time.Time, error) {
	ok, err := service.accountRepo.AccountExistsByEmail(ctx, email)
	if err != nil {
		service.log.Error().Err(err).Send()
		return "", time.Now(), err
	}

	if !ok {
		return "", time.Now(), domain.ErrNoAccount
	}

	ok, err = service.loginCodeRepo.LoginCodeExists(ctx, "web_login_code", email)
	if err != nil {
		service.log.Error().Err(err).Send()
		return "", time.Now(), err
	}

	if ok {
		return "", time.Now(), domain.ErrLoginCodeExist
	}

	code := utils.MustNewCode()
	lifetime := typedenv.Duration("LOGIN_CODE_LIFETIME")

	err = service.loginCodeRepo.AddLoginCode(ctx, "web_login_code", email, code, lifetime)
	if err == nil {
		return code, time.Now().Add(lifetime), nil
	}

	switch {
	case errors.Is(err, domain.ErrLoginCodeExist):
	default:
		service.log.Error().Err(err).Send()
	}
	return "", time.Now(), err
}

func (service *AuthService) LoginCode(ctx context.Context, email, code string) (string, string, error) {
	loginCode, err := service.loginCodeRepo.GetLoginCode(ctx, "web_login_code", email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoLoginCode):
		default:
			service.log.Error().Err(err).Send()
		}
		return "", "", err
	}

	if loginCode != code {
		return "", "", domain.ErrWrongCode
	}

	account, err := service.accountRepo.GetAccountByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoAccount):
		default:
			service.log.Error().Err(err).Send()
		}
		return "", "", err
	}

	now := time.Now()

	access, err := service.GenerateToken(map[string]interface{}{
		"iss":   "nuclear",
		"id":    account.Id,
		"email": account.Email,
		"perms": account.Perms,
		"exp":   now.Add(typedenv.Duration("ACCESS_LIFETIME")).Unix(),
		"iat":   now.Unix(),
	}, typedenv.String("ACCESS_KEY"))

	if err != nil {
		service.log.Error().Err(err).Send()
		return "", "", err
	}

	lifetime := typedenv.Duration("REFRESH_LIFETIME")
	refresh, err := service.GenerateToken(map[string]interface{}{
		"iss":   "nuclear",
		"id":    account.Id,
		"email": account.Email,
		"perms": account.Perms,
		"exp":   now.Add(lifetime).Unix(),
		"iat":   now.Unix(),
	}, typedenv.String("REFRESH_KEY"))

	if err != nil {
		service.log.Error().Err(err).Send()
		return "", "", err
	}

	err = service.tokenRepo.AddToken(ctx, "web_token", refresh, lifetime)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTokenExist):
		default:
			service.log.Error().Err(err).Send()
		}
		return "", "", err
	}

	return access, refresh, nil
}

func (service *AuthService) GenerateToken(claims map[string]interface{}, key string) (string, error) {
	var mapClaims jwt.MapClaims = claims

	unsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	signed, err := unsigned.SignedString([]byte(key))
	if err != nil {
		service.log.Error().Err(err).Send()
		return "", err
	}

	return signed, nil
}

func (service *AuthService) WhitelistToken(ctx context.Context, prefix, token string, lifetime time.Duration) error {
	err := service.tokenRepo.AddToken(ctx, prefix, token, lifetime)

	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrTokenExist):
	default:
		service.log.Error().Err(err).Send()
	}
	return err
}

func (service *AuthService) RevokeToken(ctx context.Context, prefix, token string) error {
	err := service.tokenRepo.RemoveToken(ctx, prefix, token)

	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrNoToken):
	default:
		service.log.Error().Err(err).Send()
	}
	return err
}

func (service *AuthService) GetTokenClaims(token string, key string) (map[string]interface{}, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		service.log.Error().Err(err).Send()
		return nil, err
	}

	if !parsed.Valid {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrMissingTokenClaims
	}

	return claims, nil
}

func (service *AuthService) RefreshTokens(ctx context.Context, prefix, token string) (string, string, error) {
	exists, err := service.tokenRepo.TokenExists(ctx, prefix, token)
	if err != nil || !exists {
		return "", "", domain.ErrInvalidToken
	}

	err = service.tokenRepo.RemoveToken(ctx, prefix, token)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoToken):
		default:
			service.log.Error().Err(err).Send()
		}
		return "", "", err
	}

	mapClaims, err := service.GetTokenClaims(token, typedenv.String("REFRESH_KEY"))
	if err != nil {
		return "", "", domain.ErrInvalidToken
	}

	now := time.Now()

	access, err := service.GenerateToken(map[string]interface{}{
		"iss":   "nuclear",
		"id":    mapClaims["id"],
		"email": mapClaims["email"],
		"perms": mapClaims["perms"],
		"exp":   now.Add(typedenv.Duration("ACCESS_LIFETIME")).Unix(),
		"iat":   now.Unix(),
	}, typedenv.String("ACCESS_KEY"))

	if err != nil {
		return "", "", err
	}

	lifetime := typedenv.Duration("REFRESH_LIFETIME")
	refresh, err := service.GenerateToken(map[string]interface{}{
		"iss":   "vortex",
		"id":    mapClaims["id"],
		"email": mapClaims["email"],
		"perms": mapClaims["perms"],
		"exp":   now.Add(lifetime).Unix(),
		"iat":   now.Unix(),
	}, typedenv.String("REFRESH_KEY"))

	if err != nil {
		return "", "", err
	}

	err = service.tokenRepo.AddToken(ctx, prefix, refresh, lifetime)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTokenExist):
		default:
			service.log.Error().Err(err).Send()
		}
		return "", "", err
	}

	return access, refresh, nil
}
