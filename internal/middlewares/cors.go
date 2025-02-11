package middlewares

import (
	"github.com/valyala/fasthttp"
	"strings"
)

var (
	corsAllowHeaders     = []string{"Content-Type", "Authorization"}
	corsAllowMethods     = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsAllowOrigin      = []string{"http://localhost:3000"}
	corsAllowCredentials = "true"
)

func Cors(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		setHeader := func(key string, values []string) {
			ctx.Response.Header.Set(key, strings.Join(values, ","))
		}

		setHeader(fasthttp.HeaderAccessControlAllowOrigin, corsAllowOrigin)
		setHeader(fasthttp.HeaderAccessControlAllowMethods, corsAllowMethods)
		setHeader(fasthttp.HeaderAccessControlAllowHeaders, corsAllowHeaders)

		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowCredentials, corsAllowCredentials)

		next(ctx)
	}
}
