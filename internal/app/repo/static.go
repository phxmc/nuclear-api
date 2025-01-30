package repo

import "context"

type StaticReader interface {
	GetAvatar(ctx context.Context, accountId string) ([]byte, error)
	GetBanner(ctx context.Context, accountId string) ([]byte, error)
}

type StaticWriter interface {
	SetAvatar(ctx context.Context, accountId string, avatar []byte) error
	SetBanner(ctx context.Context, accountId string, banner []byte) error

	DelAvatar(ctx context.Context, accountId string) error
	DelBanner(ctx context.Context, accountId string) error
}

type StaticReadWriter interface {
	StaticReader
	StaticWriter
}
