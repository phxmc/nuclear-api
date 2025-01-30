package builders

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/services"
)

type AuthServiceBuilder interface {
	Builder[api.AuthApi]
	AccountRepo(accountRepo repo.AccountReadWriter) AuthServiceBuilder
	LoginCodeRepo(loginCodeRepo repo.LoginCodeReadWriter) AuthServiceBuilder
	TokenRepo(tokenRepo repo.TokenReadWriter) AuthServiceBuilder
}

type authServiceBuilder struct {
	accountRepo   repo.AccountReadWriter
	loginCodeRepo repo.LoginCodeReadWriter
	tokenRepo     repo.TokenReadWriter
}

func NewAuthServiceBuilder() AuthServiceBuilder {
	return &authServiceBuilder{}
}

func (builder *authServiceBuilder) Build() api.AuthApi {
	return services.NewAuthService(
		builder.accountRepo,
		builder.loginCodeRepo,
		builder.tokenRepo,
	)
}

func (builder *authServiceBuilder) AccountRepo(accountRepo repo.AccountReadWriter) AuthServiceBuilder {
	builder.accountRepo = accountRepo
	return builder
}

func (builder *authServiceBuilder) LoginCodeRepo(loginCodeRepo repo.LoginCodeReadWriter) AuthServiceBuilder {
	builder.loginCodeRepo = loginCodeRepo
	return builder
}

func (builder *authServiceBuilder) TokenRepo(tokenRepo repo.TokenReadWriter) AuthServiceBuilder {
	builder.tokenRepo = tokenRepo
	return builder
}
