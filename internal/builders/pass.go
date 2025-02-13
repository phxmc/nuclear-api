package builders

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/services"
	"github.com/rs/zerolog"
)

type PassServiceBuilder interface {
	Builder[api.PassApi]
	PassRepo(passRepo repo.PassReadWriter) PassServiceBuilder
	Log(log *zerolog.Logger) PassServiceBuilder
}

type passServiceBuilder struct {
	passRepo repo.PassReadWriter
	log      *zerolog.Logger
}

func NewPassServiceBuilder() PassServiceBuilder {
	return &passServiceBuilder{}
}

func (builder *passServiceBuilder) Build() api.PassApi {
	return services.NewPassService(
		builder.passRepo,
		builder.log,
	)
}

func (builder *passServiceBuilder) PassRepo(passRepo repo.PassReadWriter) PassServiceBuilder {
	builder.passRepo = passRepo
	return builder
}

func (builder *passServiceBuilder) Log(log *zerolog.Logger) PassServiceBuilder {
	builder.log = log
	return builder
}
