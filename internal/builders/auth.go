package builders

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/services"
	"github.com/rs/zerolog"
)

type AuthApiBuilder interface {
	Builder[api.AuthApi]
	AccountRepo(repo.AccountReadWriter) AuthApiBuilder
	LoginCodeRepo(repo.LoginCodeReadWriter) AuthApiBuilder
	TokenRepo(repo.TokenReadWriter) AuthApiBuilder
	Log(*zerolog.Logger) AuthApiBuilder
}

type authApiBuilder struct {
	accountRepo   repo.AccountReadWriter
	loginCodeRepo repo.LoginCodeReadWriter
	tokenRepo     repo.TokenReadWriter
	log           *zerolog.Logger
}

func NewAuthServiceBuilder() AuthApiBuilder {
	return &authApiBuilder{}
}

func (builder *authApiBuilder) Build() api.AuthApi {
	return services.NewAuthService(
		builder.accountRepo,
		builder.loginCodeRepo,
		builder.tokenRepo,
		builder.log,
	)
}

func (builder *authApiBuilder) AccountRepo(accountRepo repo.AccountReadWriter) AuthApiBuilder {
	builder.accountRepo = accountRepo
	return builder
}

func (builder *authApiBuilder) LoginCodeRepo(loginCodeRepo repo.LoginCodeReadWriter) AuthApiBuilder {
	builder.loginCodeRepo = loginCodeRepo
	return builder
}

func (builder *authApiBuilder) TokenRepo(tokenRepo repo.TokenReadWriter) AuthApiBuilder {
	builder.tokenRepo = tokenRepo
	return builder
}

func (builder *authApiBuilder) Log(log *zerolog.Logger) AuthApiBuilder {
	builder.log = log
	return builder
}
