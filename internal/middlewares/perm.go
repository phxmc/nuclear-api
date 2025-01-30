package middlewares

import (
	"fmt"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

func Perm(permGroup *domain.PermGroup, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		perms, err := strconv.Atoi(fmt.Sprintf("%d", ctx.UserValue("perms")))

		if err != nil {
			utils.MustWriteString(ctx, "invalid perms", fasthttp.StatusUnauthorized)
			return
		}

		ok := true

		switch permGroup.GroupMode {
		case domain.GroupModeAll:
			ok = true

			for _, perm := range permGroup.Perms {
				if !domain.HasPerm(perms, perm) {
					ok = false
				}
			}

			if !ok {
				utils.MustWriteString(ctx, "permission denied", http.StatusForbidden)
				return
			}
			break

		case domain.GroupModeAny:
			ok = false

			for _, perm := range permGroup.Perms {
				if domain.HasPerm(perms, perm) {
					ok = true
				}
			}

			if !ok {
				utils.MustWriteString(ctx, "permission denied", http.StatusForbidden)
				return
			}
			break
		}

		handler(ctx)
	}
}
