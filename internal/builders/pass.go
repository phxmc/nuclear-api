package builders

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/services"
	"github.com/rs/zerolog"
)

type PassApiBuilder interface {
	Builder[api.PassApi]
	PassRepo(repo.PassReadWriter) PassApiBuilder
	Log(*zerolog.Logger) PassApiBuilder
}

type passApiBuilder struct {
	passRepo repo.PassReadWriter
	log      *zerolog.Logger
}

func NewPassApiBuilder() PassApiBuilder {
	return &passApiBuilder{}
}

func (builder *passApiBuilder) Build() api.PassApi {
	return services.NewPassService(
		builder.passRepo,
		builder.log,
	)
}

func (builder *passApiBuilder) PassRepo(passRepo repo.PassReadWriter) PassApiBuilder {
	builder.passRepo = passRepo
	return builder
}

func (builder *passApiBuilder) Log(log *zerolog.Logger) PassApiBuilder {
	builder.log = log
	return builder
}
