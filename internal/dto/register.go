package dto

import (
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/validator"
	"regexp"
	"time"
)

type RegisterRequest struct {
	Email string `json:"email"`
}

// Validate may return domain.ErrIncorrectEmail
func (request *RegisterRequest) Validate() error {
	ok, err := regexp.MatchString(validator.EmailRegexp, request.Email)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrIncorrectEmail
	}

	return nil
}

type RegisterResponse struct {
	Deadline time.Time `json:"deadline"`
}

type RegisterCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// Validate may return domain.ErrIncorrectEmail
func (request *RegisterCodeRequest) Validate() error {
	ok, err := regexp.MatchString(validator.EmailRegexp, request.Email)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrIncorrectEmail
	}

	return nil
}
