package builders

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/services"
	"github.com/rs/zerolog"
)

type NicknameApiBuilder interface {
	Builder[api.NicknameApi]
	NicknameRepo(repo.NicknameReadWriter) NicknameApiBuilder
	Log(*zerolog.Logger) NicknameApiBuilder
}

type nicknameApiBuilder struct {
	nicknameRepo repo.NicknameReadWriter
	log          *zerolog.Logger
}

func NewNicknameApiBuilder() NicknameApiBuilder {
	return &nicknameApiBuilder{}
}

func (builder *nicknameApiBuilder) Build() api.NicknameApi {
	return services.NewNicknameService(
		builder.nicknameRepo,
		builder.log,
	)
}

func (builder *nicknameApiBuilder) NicknameRepo(passRepo repo.NicknameReadWriter) NicknameApiBuilder {
	builder.nicknameRepo = passRepo
	return builder
}

func (builder *nicknameApiBuilder) Log(log *zerolog.Logger) NicknameApiBuilder {
	builder.log = log
	return builder
}
