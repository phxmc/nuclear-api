package repo

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
)

type PassReader interface {
	GetPassById(ctx context.Context, id string) (*domain.Pass, error)
	GetPassByAccountId(ctx context.Context, accountId string) (*domain.Pass, error)
}

type PassWriter interface {
	AddPass(ctx context.Context, pass *domain.Pass) error
}

type PassReadWriter interface {
	PassReader
	PassWriter
}
