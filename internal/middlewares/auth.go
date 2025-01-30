package middlewares

import (
	"fmt"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/orewaee/typedenv"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Auth(authApi api.AuthApi, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		header := string(ctx.Request.Header.Peek("Authorization"))
		if header == "" {
			utils.MustWriteString(ctx, "missing authorization header", fasthttp.StatusUnauthorized)
			return
		}

		accessToken := strings.TrimPrefix(header, "Bearer ")
		if accessToken == "" {
			utils.MustWriteString(ctx, "missing token", fasthttp.StatusUnauthorized)
			return
		}

		if accessToken == "" {
			utils.MustWriteString(ctx, "missing token", http.StatusUnauthorized)
			return
		}

		fmt.Println(accessToken)

		claims, err := authApi.GetTokenClaims(accessToken, typedenv.String("ACCESS_KEY"))
		if err != nil {
			log.Println(err)
			utils.MustWriteString(ctx, "invalid token", http.StatusUnauthorized)
			return
		}

		permsClaim, ok := claims["perms"]
		if !ok {
			utils.MustWriteString(ctx, "invalid token", http.StatusUnauthorized)
			return
		}

		perms, err := strconv.Atoi(fmt.Sprintf("%v", permsClaim))
		if err != nil {
			utils.MustWriteString(ctx, "invalid token", http.StatusUnauthorized)
			return
		}

		emailClaim, ok := claims["email"]
		if !ok {
			utils.MustWriteString(ctx, "invalid token", http.StatusUnauthorized)
			return
		}

		email := fmt.Sprintf("%s", emailClaim)

		ctx.SetUserValue("email", email)
		ctx.SetUserValue("perms", perms)

		handler(ctx)
	}
}
