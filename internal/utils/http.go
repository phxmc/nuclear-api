package utils

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

// MustWriteBytes writes a slice of bytes to the writer and sets the status code
func MustWriteBytes(ctx *fasthttp.RequestCtx, data []byte, code int) {
	ctx.SetStatusCode(code)
	if _, err := ctx.Write(data); err != nil {
		panic(err)
	}
}

// MustWriteString writes the string to the writer and sets the status code
func MustWriteString(ctx *fasthttp.RequestCtx, data string, code int) {
	MustWriteBytes(ctx, []byte(data), code)
}

// MustWriteJson writes data in json format to the writer and sets the status code
func MustWriteJson(ctx *fasthttp.RequestCtx, data interface{}, code int) {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	ctx.Response.Header.Set(fasthttp.HeaderContentType, "application/json")
	MustWriteBytes(ctx, bytes, code)
}

// MustReadJson reads data from the request body in json format
func MustReadJson[T interface{}](ctx *fasthttp.RequestCtx) *T {
	bytes := ctx.PostBody()

	data := new(T)
	if err := json.Unmarshal(bytes, data); err != nil {
		panic(err)
	}

	return data
}
