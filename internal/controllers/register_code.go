package controllers

import (
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
)

func (controller *RestController) registerCode(ctx *fasthttp.RequestCtx) {
	data := utils.MustReadJson[dto.RegisterCodeRequest](ctx)

	if err := data.Validate(); err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrIncorrectEmail):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusBadRequest)
			return
		default:
			controller.log.Error().Err(err).Send()

			response.Message = domain.ErrUnexpectedError.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
			return
		}
	}

	account, err := controller.accountApi.SaveTempAccount(ctx, data.Email, data.Code)
	if err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrWrongCode) || errors.Is(err, domain.ErrTempAccountNotExist):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusNotFound)
			return
		default:
			controller.log.Error().Err(err).Send()

			response.Message = domain.ErrUnexpectedError.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
			return
		}
	}

	// success message
	go controller.emailApi.Send(ctx, data.Email, "Welcome", "success")

	controller.log.Info().
		Str("id", account.Id).
		Str("email", account.Email).
		Msg("new account registered")

	accountDto := &dto.Account{
		Id:    account.Id,
		Email: account.Email,
		Perms: account.Perms,
	}

	utils.MustWriteJson(ctx, accountDto, fasthttp.StatusCreated)
}
