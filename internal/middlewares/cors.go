package middlewares

import (
	"github.com/orewaee/typedenv"
	"github.com/valyala/fasthttp"
)

var (
	corsAllowOrigin      = typedenv.String("CORS_ORIGINS", "http://localhost:3000")
	corsAllowMethods     = typedenv.String("CORS_METHODS", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	corsAllowHeaders     = typedenv.String("CORS_HEADERS", "Content-Type, Authorization")
	corsAllowCredentials = typedenv.String("CORS_CREDENTIALS", "true")
)

func Cors(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		setHeader := func(key string, value string) {
			ctx.Response.Header.Set(key, value)
		}

		setHeader(fasthttp.HeaderAccessControlAllowOrigin, corsAllowOrigin)
		setHeader(fasthttp.HeaderAccessControlAllowMethods, corsAllowMethods)
		setHeader(fasthttp.HeaderAccessControlAllowHeaders, corsAllowHeaders)
		setHeader(fasthttp.HeaderAccessControlAllowCredentials, corsAllowCredentials)

		next(ctx)
	}
}
