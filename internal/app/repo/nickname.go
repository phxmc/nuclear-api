package repo

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/domain"
)

type NicknameReader interface {
	// GetNicknameByAccountId returns the active nickname for the specified account.
	//
	// May return domain.ErrNoAccount, domain.ErrNoNickname.
	GetNicknameByAccountId(ctx context.Context, accountId string) (*domain.Nickname, error)

	// GetNicknameHistoryByAccountId returns the nickname history for the specified account in chronological order.
	// All nicknames in history are inactive.
	//
	// May return domain.ErrNoAccount.
	GetNicknameHistoryByAccountId(ctx context.Context, accountId string) ([]*domain.Nickname, error)

	// NicknameExists returns a bool value indicating the existence of the specified nickname.
	NicknameExists(ctx context.Context, nickname string) (bool, error)
}

type NicknameWriter interface {
	// SetNickname creates a new nickname and sets it to the specified account.
	// If an active nickname was associated with the account, it will be marked as inactive.
	//
	// May return domain.ErrNoAccount.
	SetNickname(ctx context.Context, accountId string, nickname *domain.Nickname) error
}

type NicknameReadWriter interface {
	NicknameReader
	NicknameWriter
}
