package utils

import (
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

func ValidatePass(pass *domain.Pass) error {
	if pass == nil {
		return domain.ErrInvalidPass
	}

	now := time.Now()
	from := pass.From
	to := pass.To

	low := false
	if from == nil {
		low = true
	} else if now.After(*from) {
		low = true
	}

	high := false
	if to == nil {
		high = true
	} else if now.Before(*to) {
		high = true
	}

	if !low || !high {
		return domain.ErrInvalidPass
	}

	return nil
}
