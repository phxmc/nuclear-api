package middlewares

import "github.com/valyala/fasthttp"

type Middleware interface {
	Use(handler fasthttp.RequestHandler) fasthttp.RequestHandler
}
