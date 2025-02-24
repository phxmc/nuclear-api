package utils

import (
	"encoding/json"
	"github.com/orewaee/nuclear-api/internal/app/domain"
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

	if len(bytes) == 0 {
		return nil
	}

	data := new(T)
	if err := json.Unmarshal(bytes, data); err != nil {
		panic(err)
	}
	return data
}

func ExtractPerms(ctx *fasthttp.RequestCtx) (int, error) {
	value := ctx.UserValue("perms")
	if value == nil {
		return 0, domain.ErrNoPerms
	}

	perms, ok := value.(int)
	if !ok {
		return 0, domain.ErrNoPerms
	}

	return perms, nil
}

func ExtractId(ctx *fasthttp.RequestCtx) (string, error) {
	value := ctx.Value("id")
	if value == nil {
		return "", domain.ErrNoId
	}

	id, ok := value.(string)
	if !ok {
		return "", domain.ErrNoId
	}

	return id, nil
}
