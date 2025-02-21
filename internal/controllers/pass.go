package controllers

import (
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
	"net/http"
)

func (controller *RestController) setPass(ctx *fasthttp.RequestCtx) {
	data := utils.MustReadJson[dto.PassRequest](ctx)
	if data == nil {
		response := &dto.Error{Message: "missing request body"}
		utils.MustWriteJson(ctx, response, fasthttp.StatusBadRequest)
		return
	}

	pass, err := controller.passApi.SetPass(ctx, data.AccountId, data.From, data.To)

	if err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrNoAccount):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, http.StatusNotFound)
			return
		case errors.Is(err, domain.ErrPassExist):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, http.StatusConflict)
			return
		default:
			controller.log.Error().Err(err).Send()
			response.Message = domain.ErrUnexpected.Error()
			utils.MustWriteJson(ctx, response, http.StatusInternalServerError)
			return
		}
	}

	response := &dto.Pass{
		Id:        pass.Id,
		From:      pass.From,
		To:        pass.To,
		CreatedAt: pass.CreatedAt,
	}

	utils.MustWriteJson(ctx, response, fasthttp.StatusCreated)
}
