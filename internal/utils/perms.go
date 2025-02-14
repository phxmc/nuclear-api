package utils

import (
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/valyala/fasthttp"
)

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
