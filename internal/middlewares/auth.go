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

func Auth(authApi api.AuthApi, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
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

		claims, err := authApi.GetTokenClaims(accessToken, typedenv.String("ACCESS_KEY"))
		if err != nil {
			response := &dto.Error{Message: "invalid token"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

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

		emailClaim, ok := claims["email"]
		if !ok {
			response := &dto.Error{Message: "invalid token"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		email := fmt.Sprintf("%s", emailClaim)

		ctx.SetUserValue("email", email)
		ctx.SetUserValue("perms", perms)

		handler(ctx)
	}
}
