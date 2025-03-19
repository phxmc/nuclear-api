package controllers

import (
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
)

func (controller *RestController) join(ctx *fasthttp.RequestCtx) {
	data := utils.MustReadJson[dto.JoinRequest](ctx)
	if data == nil {
		response := &dto.Error{Message: "missing request body"}
		utils.MustWriteJson(ctx, response, fasthttp.StatusBadRequest)
		return
	}

	accountExists, err := controller.accountApi.AccountExistsById(ctx, data.AccountId)
	if err != nil {
		response := &dto.Error{Message: domain.ErrUnexpected.Error()}
		controller.log.Error().Err(err).Send()
		utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
		return
	}

	if !accountExists {
		response := &dto.Error{Message: domain.ErrNoAccount.Error()}
		utils.MustWriteJson(ctx, response, fasthttp.StatusForbidden)
		return
	}

	_, err = controller.passApi.GetPassByAccountId(ctx, data.AccountId)
	if err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrNoAccount):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusForbidden)
			return
		case errors.Is(err, domain.ErrNoPass), errors.Is(err, domain.ErrInvalidPass):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusForbidden)
			return
		default:
			controller.log.Error().Err(err).Send()
			response.Message = domain.ErrUnexpected.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
			return
		}
	}

	nickname, err := controller.nicknameApi.GetNicknameByAccountId(ctx, data.AccountId)
	if err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrNoAccount):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusForbidden)
			return
		case errors.Is(err, domain.ErrNoNickname):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusForbidden)
			return
		default:
			controller.log.Error().Err(err).Send()
			response.Message = domain.ErrUnexpected.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
			return
		}
	}

	if nickname.Value != data.Nickname {
		response := &dto.Error{Message: "invalid nickname"}
		utils.MustWriteJson(ctx, response, fasthttp.StatusForbidden)
		return
	}

	response := &dto.JoinResponse{CanJoin: true}
	utils.MustWriteJson(ctx, response, fasthttp.StatusOK)
}
