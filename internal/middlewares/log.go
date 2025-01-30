package middlewares

import (
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"time"
)

func Log(log *zerolog.Logger, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		now := time.Now()
		handler(ctx)
		elapsed := time.Since(now)

		log.Debug().
			Str("method", string(ctx.Method())).
			Str("path", string(ctx.Path())).
			Str("elapsed", elapsed.String()).
			Send()
	}
}
