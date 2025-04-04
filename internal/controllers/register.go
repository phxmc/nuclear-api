package controllers

import (
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/orewaee/typedenv"
	"github.com/valyala/fasthttp"
	"time"
)

func (controller *RestController) register(ctx *fasthttp.RequestCtx) {
	data := utils.MustReadJson[dto.RegisterRequest](ctx)
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

	lifetime := typedenv.Duration("TEMP_ACCOUNT_LIFETIME", time.Minute)
	tempAccount, deadline, err := controller.accountApi.AddTempAccount(ctx, data.Email, lifetime)

	if err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrTempAccountExist) || errors.Is(err, domain.ErrAccountExist):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusConflict)
			return
		default:
			controller.log.Error().Err(err).Send()

			response.Message = domain.ErrUnexpected.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
			return
		}
	}

	device := string(ctx.UserAgent()) + " " + ctx.RemoteAddr().String()
	go controller.emailApi.SendRegisterEmail(ctx, data.Email, device, time.Now().Format(time.RFC822), tempAccount.Code)

	response := &dto.RegisterResponse{Deadline: deadline}
	utils.MustWriteJson(ctx, response, fasthttp.StatusCreated)
}
