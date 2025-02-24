package controllers

import (
	fastrouter "github.com/fasthttp/router"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/middlewares"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

type RestController struct {
	addr       string
	authApi    api.AuthApi
	accountApi api.AccountApi
	emailApi   api.EmailApi
	staticApi  api.StaticApi
	passApi    api.PassApi
	log        *zerolog.Logger
}

func NewRestController(
	addr string,
	authApi api.AuthApi,
	accountApi api.AccountApi,
	emailApi api.EmailApi,
	staticApi api.StaticApi,
	passApi api.PassApi,
	log *zerolog.Logger) *RestController {
	return &RestController{
		addr,
		authApi,
		accountApi,
		emailApi,
		staticApi,
		passApi,
		log,
	}
}

func (controller *RestController) Run() error {
	router := fastrouter.New()

	optionsHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	router.OPTIONS("/*", optionsHandler)

	router.GET("/ping", func(ctx *fasthttp.RequestCtx) {
		utils.MustWriteString(ctx, "pong", fasthttp.StatusOK)
	})

	v1 := router.Group("/v1")

	v1.POST("/register", controller.register)
	v1.POST("/register/code", controller.registerCode)
	v1.POST("/login", controller.login)
	v1.POST("/login/code", controller.loginCode)
	v1.POST("/refresh", controller.refresh)

	v1.GET("/me", middlewares.Auth(controller.authApi, controller.me))

	passPerm := middlewares.NewPermMiddleware(&domain.PermGroup{
		Perms:     []int{domain.PermSuper},
		GroupMode: domain.GroupModeAll,
	})
	v1.GET("/pass", middlewares.Auth(controller.authApi, controller.getPass))
	v1.GET("/pass/history", middlewares.Auth(controller.authApi, controller.getPassHistory))
	v1.POST("/pass", middlewares.Auth(controller.authApi, passPerm.Use(controller.setPass)))

	v1.GET("/avatar/{account_id}", controller.getAvatar)
	v1.POST("/avatar", middlewares.Auth(controller.authApi, controller.setAvatar))

	v1.GET("/banner/{account_id}", controller.getBanner)
	v1.POST("/banner", middlewares.Auth(controller.authApi, controller.setBanner))

	controller.log.Info().Msgf("running app on addr %s", controller.addr)
	return fasthttp.ListenAndServe(controller.addr, middlewares.Cors(middlewares.Log(controller.log, router.Handler)))
}
