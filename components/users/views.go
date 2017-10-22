package users

import (
	"donategold.me/components"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Get("/", usersGet)
	api.Get("/user", userGet)
	api.Get("/user/providers", getProviders)
	api.Post("/logout", logout)
	api.Delete("/providers/:id", deleteProvider)
}

func usersGet(c *iris.Context) {
	c.JSON(200, getAllUsers())
}

func userGet(ctx *iris.Context) {
	user := ctx.Get("token").(AccessToken)
	timestamp, tokenCentrgo := components.GetTokenForCentrifugo(user.User.ID)
	userProviders := GetProvidersByUser(user.User.ID)
	centrifugo := map[string]string{"timestamp": timestamp, "token": tokenCentrgo}
	ctx.JSON(200, map[string]interface{}{"user": user, "centrifugo": centrifugo, "providers": userProviders})
}

func logout(c *iris.Context) {
	user := c.Get("token").(AccessToken)
	err := user.Delete()
	if err != nil {
		c.JSON(404, err)
	}
	c.RemoveCookie("token")
	c.JSON(200, "")
}

func deleteProvider(c *iris.Context) {
	user := c.Get("token").(AccessToken)
	id, err := c.ParamInt("id")
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	err = deleteProviderByIdAndUserId(user.User.ID, uint(id))
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	c.JSON(200, Provider{ID: uint(id)})
}

func getProviders(c *iris.Context) {
	user := c.Get("token").(AccessToken)
	providers := GetProvidersByUser(user.User.ID)
	c.JSON(200, providers)
}
