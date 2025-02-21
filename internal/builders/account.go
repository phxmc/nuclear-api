package builders

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/services"
	"github.com/rs/zerolog"
)

type AccountApiBuilder interface {
	Builder[api.AccountApi]
	AccountRepo(repo.AccountReadWriter) AccountApiBuilder
	TempAccountRepo(repo.TempAccountReadWriter) AccountApiBuilder
	Log(*zerolog.Logger) AccountApiBuilder
}

type accountApiBuilder struct {
	accountRepo     repo.AccountReadWriter
	tempAccountRepo repo.TempAccountReadWriter
	log             *zerolog.Logger
}

func NewAccountApiBuilder() AccountApiBuilder {
	return &accountApiBuilder{}
}

func (builder *accountApiBuilder) Build() api.AccountApi {
	return services.NewAccountService(
		builder.accountRepo,
		builder.tempAccountRepo,
		builder.log,
	)
}

func (builder *accountApiBuilder) AccountRepo(accountRepo repo.AccountReadWriter) AccountApiBuilder {
	builder.accountRepo = accountRepo
	return builder
}

func (builder *accountApiBuilder) TempAccountRepo(tempAccountRepo repo.TempAccountReadWriter) AccountApiBuilder {
	builder.tempAccountRepo = tempAccountRepo
	return builder
}

func (builder *accountApiBuilder) Log(log *zerolog.Logger) AccountApiBuilder {
	builder.log = log
	return builder
}
