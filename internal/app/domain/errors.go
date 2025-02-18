package domain

import "errors"

var (
	ErrTempAccountAlreadyExists = errors.New("account already exists")

	ErrTempAccountExist = errors.New("temp account already exists")
	ErrNoTempAccount    = errors.New("temp account does not exist")

	ErrAccountExist = errors.New("account already exists")
	ErrNoAccount    = errors.New("account does not exist")

	ErrPassExist = errors.New("pass already exists")
	ErrNoPass    = errors.New("pass does not exist")

	ErrInvalidPass = errors.New("invalid pass")

	ErrTokenExist = errors.New("token already exists")
	ErrNoToken    = errors.New("token does not exist")

	ErrLoginCodeExist = errors.New("login code already exists")
	ErrNoLoginCode    = errors.New("login code does not exist")

	ErrWrongCode = errors.New("wrong code")

	ErrInvalidToken       = errors.New("invalid token")
	ErrMissingTokenClaims = errors.New("missing token claims")

	ErrAvatarNotExist = errors.New("avatar does not exist")
	ErrInvalidAvatar  = errors.New("invalid avatar")

	ErrBannerNotExist = errors.New("banner does not exist")
	ErrInvalidBanner  = errors.New("invalid banner")

	ErrIncorrectEmail   = errors.New("incorrect email")
	ErrTempCodeNotFound = errors.New("temp code not found")

	ErrUnexpected = errors.New("unexpected error")

	ErrNoPerms = errors.New("permissions does not exist")
)
