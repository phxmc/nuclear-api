package driving

import "context"

type Bot interface {
	SendMessage(ctx context.Context, markdown string) error
	Run(ctx context.Context, token string) error
}
