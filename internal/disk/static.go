package disk

import (
	"context"
	"errors"
	"fmt"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"os"
)

type StaticRepo struct{}

func NewStaticRepo() repo.StaticReadWriter {
	return &StaticRepo{}
}

func (repo *StaticRepo) GetAvatar(ctx context.Context, accountId string) ([]byte, error) {
	path := fmt.Sprintf("./static/avatar-%s", accountId)

	info, err := os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil, domain.ErrAvatarNotExist
	}

	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, domain.ErrInvalidAvatar
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (repo *StaticRepo) GetBanner(ctx context.Context, accountId string) ([]byte, error) {
	path := fmt.Sprintf("./static/banner-%s", accountId)

	info, err := os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil, domain.ErrBannerNotExist
	}

	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, domain.ErrInvalidBanner
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (repo *StaticRepo) SetAvatar(ctx context.Context, accountId string, avatar []byte) error {
	path := fmt.Sprintf("./static/avatar-%s", accountId)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := file.Write(avatar); err != nil {
		return err
	}

	return file.Close()
}

func (repo *StaticRepo) SetBanner(ctx context.Context, accountId string, banner []byte) error {
	path := fmt.Sprintf("./static/banner-%s", accountId)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := file.Write(banner); err != nil {
		return err
	}

	return file.Close()
}

func (repo *StaticRepo) DelAvatar(ctx context.Context, accountId string) error {
	path := fmt.Sprintf("./static/avatar-%s", accountId)
	return os.Remove(path)
}

func (repo *StaticRepo) DelBanner(ctx context.Context, accountId string) error {
	path := fmt.Sprintf("./static/banner-%s", accountId)
	return os.Remove(path)
}
