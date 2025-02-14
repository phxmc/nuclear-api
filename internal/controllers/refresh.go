package controllers

import (
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
)

func (controller *RestController) refresh(ctx *fasthttp.RequestCtx) {
	data := utils.MustReadJson[dto.RefreshRequest](ctx)

	access, refresh, err := controller.authApi.RefreshToken(ctx, data.RefreshToken)
	if err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrInvalidToken):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
			return
		default:
			controller.log.Error().Err(err).Send()

			response.Message = domain.ErrUnexpected.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
			return
		}
	}

	pair := &dto.TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	utils.MustWriteJson(ctx, pair, fasthttp.StatusOK)
}
