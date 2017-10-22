package money

import (
	"fmt"

	"donategold.me/components/users"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Get("/balance", GetUserBalance)
	//api.Get("/add/:username", addBalance)
}

func GetUserBalance(c *iris.Context) {
	user := c.Get("token").(users.AccessToken)
	if !user.IsAuthenticated() {
		c.EmitError(404)
		return
	}
	balance := GetBalanceByUserId(user.User.ID)
	c.JSON(200, balance)
}

func addBalance(c *iris.Context) {
	user, exist := users.GetUserByUsername(c.Param("username"))
	if !exist {
		c.EmitError(404)
		return
	}
	balance := GetBalanceByUserId(user.ID)
	money, err := c.URLParamInt("money")
	if err != nil {
		c.EmitError(404)
		return
	}
	balance.Add(money)
	if err != nil {
		fmt.Println(err)
		c.EmitError(404)
		return
	}
	c.WriteString("ok")
}
