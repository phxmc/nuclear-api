package services

import (
	"context"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"os"
)

type StaticService struct {
	staticRepo repo.StaticReadWriter
}

func NewStaticService(staticRepo repo.StaticReadWriter) api.StaticApi {
	os.Mkdir("static", 0700)
	return &StaticService{staticRepo}
}

func (service *StaticService) GetAvatar(ctx context.Context, accountId string) ([]byte, error) {
	return service.staticRepo.GetAvatar(ctx, accountId)
}

func (service *StaticService) GetBanner(ctx context.Context, accountId string) ([]byte, error) {
	return service.staticRepo.GetBanner(ctx, accountId)
}

func (service *StaticService) SetAvatar(ctx context.Context, accountId string, avatar []byte) error {
	return service.staticRepo.SetAvatar(ctx, accountId, avatar)
}

func (service *StaticService) SetBanner(ctx context.Context, accountId string, banner []byte) error {
	return service.staticRepo.SetBanner(ctx, accountId, banner)
}

func (service *StaticService) DelAvatar(ctx context.Context, accountId string) error {
	return service.staticRepo.DelAvatar(ctx, accountId)
}

func (service *StaticService) DelBanner(ctx context.Context, accountId string) error {
	return service.staticRepo.DelBanner(ctx, accountId)
}
