package builders

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/services"
)

type AccountServiceBuilder interface {
	Builder[api.AccountApi]
	AccountRepo(accountRepo repo.AccountReadWriter) AccountServiceBuilder
	TempAccountRepo(tempAccountRepo repo.TempAccountReadWriter) AccountServiceBuilder
}

type accountServiceBuilder struct {
	accountRepo     repo.AccountReadWriter
	tempAccountRepo repo.TempAccountReadWriter
}

func NewAccountServiceBuilder() AccountServiceBuilder {
	return &accountServiceBuilder{}
}

func (builder *accountServiceBuilder) Build() api.AccountApi {
	return services.NewAccountService(
		builder.accountRepo,
		builder.tempAccountRepo,
	)
}

func (builder *accountServiceBuilder) AccountRepo(accountRepo repo.AccountReadWriter) AccountServiceBuilder {
	builder.accountRepo = accountRepo
	return builder
}

func (builder *accountServiceBuilder) TempAccountRepo(tempAccountRepo repo.TempAccountReadWriter) AccountServiceBuilder {
	builder.tempAccountRepo = tempAccountRepo
	return builder
}
