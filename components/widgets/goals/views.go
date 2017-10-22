package goals

import (
	"donategold.me/components/users"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Get("/", GetGoalsByUser)
	api.Post("/", NewGoal)
}

func GetGoalsByUser(c *iris.Context) {
	user := c.Get("token").(users.AccessToken)
	goals := GetGoalsByUserId(user.User.ID)
	c.JSON(200, goals)
}

func NewGoal(c *iris.Context) {
	var goal Goal
	if err := c.ReadJSON(&goal); err != nil {
		c.EmitError(404)
		return
	}
	user := c.Get("token").(users.AccessToken)
	goal.UserId = user.User.ID
	goal.Create()
	c.JSON(200, goal)
}
