package middlewares

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func Foo(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Println("foo middleware")
		handler(ctx)
	}
}

func Bar(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Println("bar middleware")
		handler(ctx)
	}
}
