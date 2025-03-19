package middlewares

import (
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/orewaee/typedenv"
	"github.com/valyala/fasthttp"
)

type ApiKeyMiddleware struct {
	apiKey string
}

func NewApiKeyMiddleware() *ApiKeyMiddleware {
	return &ApiKeyMiddleware{
		apiKey: typedenv.String("API_KEY", ""),
	}
}

func (middleware *ApiKeyMiddleware) Use(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if middleware.apiKey == "" {
			response := &dto.Error{Message: "invalid api key"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		apiKey := string(ctx.Request.Header.Peek("X-API-Key"))
		if apiKey == "" {
			response := &dto.Error{Message: "missing x-api-key"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		if apiKey != middleware.apiKey {
			response := &dto.Error{Message: "invalid api key"}
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		}

		handler(ctx)
	}
}
