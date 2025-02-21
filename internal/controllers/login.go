package controllers

import (
	"errors"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
	"time"
)

func (controller *RestController) login(ctx *fasthttp.RequestCtx) {
	data := utils.MustReadJson[dto.LoginRequest](ctx)
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

	code, deadline, err := controller.authApi.Login(ctx, data.Email)
	if err != nil {
		response := &dto.Error{}

		switch {
		case errors.Is(err, domain.ErrNoAccount):
			response.Message = err.Error()
			utils.MustWriteJson(ctx, response, fasthttp.StatusNotFound)
			return
		case errors.Is(err, domain.ErrLoginCodeExist):
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
	go controller.emailApi.SendLoginEmail(ctx, data.Email, device, time.Now().Format(time.RFC822), code)

	response := &dto.LoginResponse{Deadline: deadline}
	utils.MustWriteJson(ctx, response, fasthttp.StatusCreated)
}
