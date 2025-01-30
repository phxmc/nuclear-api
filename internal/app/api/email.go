package api

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
)

type EmailApi interface {
	Send(ctx context.Context, receiver, subject, text string) error
	SendRegisterMail(ctx context.Context, receiver string, tempAccount *domain.TempAccount) error
}
