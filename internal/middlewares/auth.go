package middlewares

import (
	"fmt"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/orewaee/typedenv"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type AuthMiddleware struct {
	authApi api.AuthApi
}

func NewAuthMiddleware(authApi api.AuthApi) Middleware {
	return &AuthMiddleware{
		authApi: authApi,
	}
}

func (middleware *AuthMiddleware) Use(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		header := string(ctx.Request.Header.Peek("Authorization"))
		if header == "" {
			response := &dto.Error{Message: "missing authorization header"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		accessToken := strings.TrimPrefix(header, "Bearer ")
		if accessToken == "" {
			response := &dto.Error{Message: "missing token"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		if accessToken == "" {
			response := &dto.Error{Message: "missing token"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		claims, err := middleware.authApi.GetTokenClaims(accessToken, typedenv.String("ACCESS_KEY"))
		if err != nil {
			response := &dto.Error{Message: "invalid token"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		idClaim, ok := claims["id"]
		if !ok {
			response := &dto.Error{Message: "invalid token"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		id := fmt.Sprintf("%s", idClaim)

		emailClaim, ok := claims["email"]
		if !ok {
			response := &dto.Error{Message: "invalid token"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		email := fmt.Sprintf("%s", emailClaim)

		permsClaim, ok := claims["perms"]
		if !ok {
			response := &dto.Error{Message: "invalid token"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		perms, err := strconv.Atoi(fmt.Sprintf("%v", permsClaim))
		if err != nil {
			response := &dto.Error{Message: "invalid token"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		ctx.SetUserValue("id", id)
		ctx.SetUserValue("email", email)
		ctx.SetUserValue("perms", perms)

		handler(ctx)
	}
}
