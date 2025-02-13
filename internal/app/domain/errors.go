package domain

import "errors"

var (
	ErrTempAccountAlreadyExists = errors.New("account already exists")

	ErrAccountExist = errors.New("account already exists")
	ErrNoAccount    = errors.New("account does not exist")

	ErrTempAccountExist    = errors.New("temp account already exists")
	ErrTempAccountNotExist = errors.New("temp account does not exist")

	ErrLoginCodeExist    = errors.New("login code already exists")
	ErrLoginCodeNotExist = errors.New("login code does not exist")

	ErrWrongCode = errors.New("wrong code")

	ErrInvalidToken       = errors.New("invalid token")
	ErrMissingTokenClaims = errors.New("missing token claims")

	ErrAvatarNotExist = errors.New("avatar does not exist")
	ErrInvalidAvatar  = errors.New("invalid avatar")

	ErrBannerNotExist = errors.New("banner does not exist")
	ErrInvalidBanner  = errors.New("invalid banner")

	ErrIncorrectEmail   = errors.New("incorrect email")
	ErrUnexpectedError  = errors.New("unexpected error")
	ErrTempCodeNotFound = errors.New("temp code not found")

	ErrNoPass    = errors.New("pass does not exist")
	ErrPassExist = errors.New("pass already exists")
)
