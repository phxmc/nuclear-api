package services

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
)

// todo refactor it

type ServiceBuilder[T interface{}] interface {
	Build() T
}

type AccountServiceBuilder interface {
	ServiceBuilder[api.AccountApi]
	AccountRepo(accountRepo repo.AccountReadWriter) AccountServiceBuilder
	TempAccountRepo(tempAccountRepo repo.TempAccountReadWriter) AccountServiceBuilder
}

type accountServiceBuilder struct {
	service *AccountService
}

func NewAccountServiceBuilder() AccountServiceBuilder {
	return &accountServiceBuilder{
		service: &AccountService{},
	}
}

func (builder *accountServiceBuilder) Build() api.AccountApi {
	return builder.service
}

func (builder *accountServiceBuilder) AccountRepo(accountRepo repo.AccountReadWriter) AccountServiceBuilder {
	builder.service.accountRepo = accountRepo
	return builder
}

func (builder *accountServiceBuilder) TempAccountRepo(tempAccountRepo repo.TempAccountReadWriter) AccountServiceBuilder {
	builder.service.tempAccountRepo = tempAccountRepo
	return builder
}

type AuthServiceBuilder interface {
	ServiceBuilder[api.AuthApi]
	AccountRepo(accountRepo repo.AccountReadWriter) AuthServiceBuilder
	LoginCodeRepo(loginCodeRepo repo.LoginCodeReadWriter) AuthServiceBuilder
	TokenRepo(tokenRepo repo.TokenReadWriter) AuthServiceBuilder
}

type authServiceBuilder struct {
	service *AuthService
}

func NewAuthServiceBuilder() AuthServiceBuilder {
	return &authServiceBuilder{
		service: &AuthService{},
	}
}

func (builder *authServiceBuilder) Build() api.AuthApi {
	return builder.service
}

func (builder *authServiceBuilder) AccountRepo(accountRepo repo.AccountReadWriter) AuthServiceBuilder {
	builder.service.accountRepo = accountRepo
	return builder
}

func (builder *authServiceBuilder) LoginCodeRepo(loginCodeRepo repo.LoginCodeReadWriter) AuthServiceBuilder {
	builder.service.loginCodeRepo = loginCodeRepo
	return builder
}

func (builder *authServiceBuilder) TokenRepo(tokenRepo repo.TokenReadWriter) AuthServiceBuilder {
	builder.service.tokenRepo = tokenRepo
	return builder
}
