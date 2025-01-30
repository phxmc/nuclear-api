package controllers

import (
	"fmt"
	"github.com/orewaee/nuclear-api/internal/dto"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
)

// todo refactor

func (controller *RestController) me(ctx *fasthttp.RequestCtx) {
	email := fmt.Sprintf("%s", ctx.UserValue("email"))

	account, err := controller.accountApi.GetAccountByEmail(ctx, email)
	if err != nil {
		utils.MustWriteString(ctx, err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	data := dto.Account{
		Id:    account.Id,
		Email: account.Email,
		Perms: account.Perms,
	}

	utils.MustWriteJson(ctx, data, fasthttp.StatusOK)
}
