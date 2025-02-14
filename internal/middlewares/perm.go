package middlewares

import (
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
	"net/http"
)

type PermMiddleware struct {
	permGroup *domain.PermGroup
}

func NewPermMiddleware(permGroup *domain.PermGroup) *PermMiddleware {
	return &PermMiddleware{
		permGroup: permGroup,
	}
}

func (middleware *PermMiddleware) Use(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		perms, err := utils.ExtractPerms(ctx)

		if err != nil {
			response := &dto.Error{}

			switch {
			case errors.Is(err, domain.ErrNoPerms):
				response.Message = err.Error()
				utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
				return
			default:
				response.Message = domain.ErrUnexpected.Error()
				utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
				return
			}
		}

		ok := true
		response := &dto.Error{}

		switch middleware.permGroup.GroupMode {
		case domain.GroupModeAll:
			ok = true

			for _, perm := range middleware.permGroup.Perms {
				if !domain.HasPerm(perms, perm) {
					ok = false
				}
			}

			if !ok {
				response.Message = "permission denied"
				utils.MustWriteJson(ctx, response, http.StatusForbidden)
				return
			}
			break

		case domain.GroupModeAny:
			ok = false

			for _, perm := range middleware.permGroup.Perms {
				if domain.HasPerm(perms, perm) {
					ok = true
				}
			}

			if !ok {
				response.Message = "permission denied"
				utils.MustWriteJson(ctx, response, http.StatusForbidden)
				return
			}
			break
		}

		handler(ctx)
	}
}
