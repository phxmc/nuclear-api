package api

import (
	"context"
)

type EmailApi interface {
	Send(ctx context.Context, receiver, subject, text string) error

	SendLoginEmail(ctx context.Context, receiver, device, datetime, code string) error
	SendRegisterEmail(ctx context.Context, receiver, device, datetime, code string) error
}
