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

	if err := data.Validate(); err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrIncorrectEmail):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusBadRequest)
			return
		default:
			response.Message = domain.ErrUnexpectedError.Error()
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
			response.Message = domain.ErrUnexpectedError.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusInternalServerError)
			return
		}
	}

	go controller.emailApi.SendRegisterMail(ctx, data.Email, tempAccount)

	response := &dto.RegisterResponse{Deadline: deadline}
	utils.MustWriteJson(ctx, response, fasthttp.StatusCreated)
}
