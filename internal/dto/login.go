package dto

import (
	"regexp"
	"time"

	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/validator"
)

type LoginMethod string

const (
	LoginMethodEmail    LoginMethod = "email"
	LoginMethodTelegram LoginMethod = "telegram"
)

type LoginRequest struct {
	Email  string      `json:"email"`
	Method LoginMethod `json:"method"`
}

// Validate may return domain.ErrIncorrectEmail
func (request *LoginRequest) Validate() error {
	ok, err := regexp.MatchString(validator.EmailRegexp, request.Email)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrIncorrectEmail
	}

	return nil
}

type LoginResponse struct {
	Deadline time.Time `json:"deadline"`
}

type LoginCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// Validate may return domain.ErrIncorrectEmail
func (request *LoginCodeRequest) Validate() error {
	ok, err := regexp.MatchString(validator.EmailRegexp, request.Email)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrIncorrectEmail
	}

	return nil
}
