package server

import (
	"donategold.me/components/auth"
	"donategold.me/components/billing"
	"donategold.me/components/money"
	"donategold.me/components/public"
	"donategold.me/components/uploader"
	"donategold.me/components/users"
	"donategold.me/components/widgets"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/middleware/logger"
)

func Get() *iris.Framework {
	server := iris.New()
	server.Adapt(iris.DevLogger())
	server.Adapt(httprouter.New())
	server.Use(logger.New())
	server.UseFunc(users.GetMiddleware(false))
	server.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		ctx.Log(iris.DevMode, "%s %d", ctx.Path(), 404)
		ctx.JSON(404, "page not found")
	})
	api := server.Party("api")
	users.Concat(api.Party("/users", users.GetMiddleware(true)))
	auth.Concat(api.Party("/auth"))
	money.Concat(api.Party("/money"))
	widgets.Concat(api.Party("/widgets", users.GetMiddleware(true)))
	uploader.Concat(api.Party("/upload", users.GetMiddleware(true)))
	billing.Concat(api.Party("/billing"))
	public.Concat(api.Party("/public"))
	return server
}
