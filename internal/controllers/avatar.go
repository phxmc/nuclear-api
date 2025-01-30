package controllers

import (
	"fmt"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/valyala/fasthttp"
)

func (controller *RestController) getAvatar(ctx *fasthttp.RequestCtx) {
	accountId := fmt.Sprintf("%s", ctx.UserValue("account_id"))

	avatar, err := controller.staticApi.GetAvatar(ctx, accountId)
	if err != nil {
		utils.MustWriteString(ctx, err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.Set("Content-Type", "image/png")
	utils.MustWriteBytes(ctx, avatar, fasthttp.StatusOK)
}

func (controller *RestController) setAvatar(ctx *fasthttp.RequestCtx) {
	email := fmt.Sprintf("%s", ctx.UserValue("email"))

	account, err := controller.accountApi.GetAccountByEmail(ctx, email)
	if err != nil {
		utils.MustWriteString(ctx, err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	contentType := string(ctx.Request.Header.Peek(fasthttp.HeaderContentType))
	if contentType != "image/png" {
		utils.MustWriteString(ctx, "invalid content-type", fasthttp.StatusUnsupportedMediaType)
		return
	}

	err = controller.staticApi.SetAvatar(ctx, account.Id, ctx.PostBody())
	if err != nil {
		utils.MustWriteString(ctx, err.Error(), fasthttp.StatusInternalServerError)
		return
	}
}
