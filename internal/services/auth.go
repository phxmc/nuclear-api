package services

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/orewaee/typedenv"
	"time"
)

type AuthService struct {
	accountRepo   repo.AccountReadWriter
	loginCodeRepo repo.LoginCodeReadWriter
	tokenRepo     repo.TokenReadWriter
}

func (service *AuthService) Login(ctx context.Context, email string) (string, time.Time, error) {
	ok, err := service.accountRepo.AccountExistsByEmail(ctx, email)
	if err != nil {
		return "", time.Now(), err
	}

	if !ok {
		return "", time.Now(), domain.ErrAccountNotExist
	}

	ok, err = service.loginCodeRepo.LoginCodeExists(ctx, email)
	if err != nil {
		return "", time.Now(), err
	}

	if ok {
		return "", time.Now(), domain.ErrLoginCodeExist
	}

	code := utils.NewCode()
	lifetime := typedenv.Duration("LOGIN_CODE_LIFETIME")
	if err := service.loginCodeRepo.AddLoginCode(ctx, email, code, lifetime); err != nil {
		return "", time.Now(), err
	}

	return code, time.Now().Add(lifetime), nil
}

func (service *AuthService) LoginCode(ctx context.Context, email, code string) (string, string, error) {
	loginCode, err := service.loginCodeRepo.GetLoginCode(ctx, email)
	if err != nil {
		return "", "", err
	}

	if loginCode != code {
		return "", "", domain.ErrWrongCode
	}

	account, err := service.accountRepo.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	now := time.Now()

	access, err := service.CreateToken(map[string]interface{}{
		"iss":   "nuclear",
		"email": account.Email,
		"perms": account.Perms,
		"exp":   now.Add(typedenv.Duration("ACCESS_LIFETIME")).Unix(),
		"iat":   now.Unix(),
	}, typedenv.String("ACCESS_KEY"))

	if err != nil {
		return "", "", err
	}

	lifetime := typedenv.Duration("REFRESH_LIFETIME")
	refresh, err := service.CreateToken(map[string]interface{}{
		"iss":   "nuclear",
		"email": account.Email,
		"perms": account.Perms,
		"exp":   now.Add(lifetime).Unix(),
		"iat":   now.Unix(),
	}, typedenv.String("REFRESH_KEY"))

	if err != nil {
		return "", "", err
	}

	if err := service.tokenRepo.AddToken(ctx, refresh, lifetime); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (service *AuthService) CreateToken(claims map[string]interface{}, key string) (string, error) {
	var mapClaims jwt.MapClaims = claims

	unsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	signed, err := unsigned.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (service *AuthService) WhitelistToken(ctx context.Context, refreshToken string, lifetime time.Duration) error {
	return service.tokenRepo.AddToken(ctx, refreshToken, lifetime)
}

func (service *AuthService) GetTokenClaims(token string, key string) (map[string]interface{}, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	})

	if err != nil {
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

func (service *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	exists, err := service.tokenRepo.TokenExists(ctx, refreshToken)

	if err != nil || !exists {
		return "", "", domain.ErrInvalidToken
	}

	if err := service.tokenRepo.RemoveToken(ctx, refreshToken); err != nil {
		return "", "", err
	}

	mapClaims, err := service.GetTokenClaims(refreshToken, typedenv.String("REFRESH_KEY"))

	if err != nil {
		return "", "", domain.ErrInvalidToken
	}

	now := time.Now()

	access, err := service.CreateToken(map[string]interface{}{
		"iss":   "nuclear",
		"email": mapClaims["email"],
		"perms": mapClaims["perms"],
		"exp":   now.Add(typedenv.Duration("ACCESS_LIFETIME")).Unix(),
		"iat":   now.Unix(),
	}, typedenv.String("ACCESS_KEY"))

	if err != nil {
		return "", "", err
	}

	lifetime := typedenv.Duration("REFRESH_LIFETIME")
	refresh, err := service.CreateToken(map[string]interface{}{
		"iss":   "vortex",
		"email": mapClaims["email"],
		"perms": mapClaims["perms"],
		"exp":   now.Add(lifetime).Unix(),
		"iat":   now.Unix(),
	}, typedenv.String("REFRESH_KEY"))

	if err != nil {
		return "", "", err
	}

	if err := service.tokenRepo.AddToken(ctx, refresh, lifetime); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
