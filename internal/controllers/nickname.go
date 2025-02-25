package controllers

import (
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
	"net/http"
)

func (controller *RestController) getNickname(ctx *fasthttp.RequestCtx) {
	id, err := utils.ExtractId(ctx)
	if err != nil {
		response := &dto.Error{Message: err.Error()}
		utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
		return
	}

	nickname, err := controller.nicknameApi.GetNicknameByAccountId(ctx, id)
	if err != nil {
		response := &dto.Error{}
		code := http.StatusInternalServerError

		switch {
		case errors.Is(err, domain.ErrNoAccount):
			response.Message = err.Error()
			code = http.StatusUnauthorized
		case errors.Is(err, domain.ErrNoNickname):
			response.Message = err.Error()
			code = http.StatusNotFound
		default:
			response.Message = domain.ErrUnexpected.Error()
			controller.log.Error().Err(err).Send()
		}

		utils.MustWriteJson(ctx, response, code)
		return
	}

	response := &dto.Nickname{
		Value:     nickname.Value,
		CreatedAt: nickname.CreatedAt,
	}

	utils.MustWriteJson(ctx, response, fasthttp.StatusOK)
}

func (controller *RestController) getNicknameHistory(ctx *fasthttp.RequestCtx) {
	id, err := utils.ExtractId(ctx)
	if err != nil {
		response := &dto.Error{Message: err.Error()}
		utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
		return
	}

	nicknames, err := controller.nicknameApi.GetNicknameHistoryByAccountId(ctx, id)
	if err != nil {
		response := &dto.Error{}
		code := http.StatusInternalServerError

		switch {
		case errors.Is(err, domain.ErrNoAccount):
			response.Message = err.Error()
			code = http.StatusUnauthorized
		default:
			controller.log.Error().Err(err).Send()
			response.Message = domain.ErrUnexpected.Error()
		}

		utils.MustWriteJson(ctx, response, code)
		return
	}

	response := make([]*dto.Nickname, len(nicknames))
	for i, nickname := range nicknames {
		response[i] = &dto.Nickname{
			Value:     nickname.Value,
			CreatedAt: nickname.CreatedAt,
		}
	}

	utils.MustWriteJson(ctx, response, fasthttp.StatusOK)
}

func (controller *RestController) setNickname(ctx *fasthttp.RequestCtx) {
	id, err := utils.ExtractId(ctx)
	if err != nil {
		response := &dto.Error{Message: err.Error()}
		utils.MustWriteJson(ctx, response, fasthttp.StatusUnauthorized)
		return
	}

	pass, err := controller.passApi.GetPassByAccountId(ctx, id)
	if err != nil {
		response := &dto.Error{}
		code := http.StatusInternalServerError

		switch {
		case errors.Is(err, domain.ErrNoAccount):
			response.Message = err.Error()
			code = http.StatusUnauthorized
		case errors.Is(err, domain.ErrNoPass):
			response.Message = err.Error()
			code = http.StatusForbidden
		default:
			controller.log.Error().Err(err).Send()
			response.Message = domain.ErrUnexpected.Error()
		}

		utils.MustWriteJson(ctx, response, code)
		return
	}

	err = utils.ValidatePass(pass)
	if err != nil {
		utils.MustWriteJson(ctx, &dto.Error{Message: err.Error()}, fasthttp.StatusForbidden)
		return
	}

	data := utils.MustReadJson[dto.NicknameRequest](ctx)
	if data == nil {
		response := &dto.Error{Message: "missing request body"}
		utils.MustWriteJson(ctx, response, fasthttp.StatusBadRequest)
		return
	}

	nickname, err := controller.nicknameApi.SetNickname(ctx, id, data.Value)
	if err != nil {
		response := &dto.Error{}
		code := http.StatusInternalServerError

		switch {
		case errors.Is(err, domain.ErrNoAccount):
			response.Message = err.Error()
			code = http.StatusNotFound
		case errors.Is(err, domain.ErrNicknameCooldown):
			response.Message = err.Error()
			code = http.StatusTooManyRequests
		default:
			controller.log.Error().Err(err).Send()
			response.Message = domain.ErrUnexpected.Error()
		}

		utils.MustWriteJson(ctx, response, code)
		return
	}

	response := &dto.Nickname{
		Value:     nickname.Value,
		CreatedAt: nickname.CreatedAt,
	}

	utils.MustWriteJson(ctx, response, fasthttp.StatusCreated)
}
