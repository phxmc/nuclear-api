package controllers

import (
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
)

func (controller *RestController) loginCode(ctx *fasthttp.RequestCtx) {
	data := utils.MustReadJson[dto.LoginCodeRequest](ctx)
	if data == nil {
		response := &dto.Error{Message: "missing request body"}
		utils.MustWriteJson(ctx, response, fasthttp.StatusBadRequest)
		return
	}

	if err := data.Validate(); err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrIncorrectEmail):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusBadRequest)
			return
		default:
			controller.log.Error().Err(err).Send()
			response.Message = domain.ErrUnexpected.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
			return
		}
	}

	access, refresh, err := controller.authApi.LoginCode(ctx, data.Email, data.Code)
	if err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrNoLoginCode) || errors.Is(err, domain.ErrWrongCode):
			response.Message = domain.ErrTempCodeNotFound.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusNotFound)
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
