package billing

import (
	"log"

	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Post("/success", successDonate)
}

func successDonate(c *iris.Context) {
	var b Bill
	c.ReadForm(&b)
	err := b.Create()
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(200, b)
}
