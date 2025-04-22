package driving

import "context"

type Bot interface {
	Run(ctx context.Context, token string) error
}
