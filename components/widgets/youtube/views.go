package youtube

import (
	"encoding/json"
	"fmt"

	"donategold.me/components/users"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Get("/", GetYoutubeWidgets)
	api.Post("/", CreateOrSave)
	api.Post("/:id/stopView", StopView)
	api.Post("/:id/viewNow", ViewNow)
	api.Post("/:id", Active)
	api.Delete("/:id", DeleteYoutubeWidget)
	api.Post("/:id/addVideo", addVideo)
	api.Post("/:id/viewedDone", YoutubeViewdVideo)
}

func CreateOrSave(c *iris.Context) {
	var yw Youtube
	c.ReadJSON(&yw)
	user_id := c.Get("token").(users.AccessToken).User.ID
	yw.UserID = user_id
	if yw.ID != 0 {
		yw.Save()
	} else {
		yw.Create()
	}
	c.JSON(200, yw)
}

func GetYoutubeWidgets(c *iris.Context) {
	user_id := c.Get("token").(users.AccessToken).User.ID
	yw := GetYoutubeWidjetsByUser(user_id)
	c.JSON(200, yw)
}

func Active(c *iris.Context) {
	user_id := c.Get("token").(users.AccessToken).User.ID
	id, err := c.ParamInt("id")
	if err != nil {
		c.EmitError(404)
		return
	}
	err = activing(user_id, id)
	if err != nil {
		c.EmitError(404)
		return
	}
	c.JSON(200, id)
}

func DeleteYoutubeWidget(c *iris.Context) {
	user_id := c.Get("token").(users.AccessToken).User.ID
	id, err := c.ParamInt("id")
	if err != nil {
		c.EmitError(404)
		return
	}
	err = deleteWidgetByUserAndId(user_id, id)
	if err != nil {
		c.EmitError(404)
		return
	}
	c.JSON(200, "")
}

func addVideo(c *iris.Context) {
	var data map[string]string
	c.ReadJSON(&data)
	id, err := c.ParamInt("id")
	if err != nil {
		c.EmitError(404)
		return
	}
	userId := c.Get("token").(users.AccessToken).User.ID
	w, err := GetYoutubeWidgetByUserAndId(userId, uint(id))
	if err != nil {
		c.EmitError(404)
		return
	}
	data["byUser"] = "true"
	err = w.AddVideo(data)
	fmt.Println(3, err)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
}

func YoutubeViewdVideo(c *iris.Context) {
	id, err := c.ParamInt("id")
	if err != nil {
		c.EmitError(404)
		return
	}
	userId := c.Get("token").(users.AccessToken).User.ID
	w, err := GetYoutubeWidgetByUserAndId(userId, uint(id))
	if err != nil {
		c.EmitError(404)
		return
	}
	var yd YoutubeDonate
	c.ReadJSON(&yd)
	err = w.ViewedVideo(yd)
	if err != nil {
		c.JSON(402, err)
		return
	}
}

func StopView(c *iris.Context) {
	userId := c.Get("token").(users.AccessToken).User.ID
	var yd YoutubeDonate
	c.ReadJSON(&yd)
	y, _ := json.Marshal(yd)
	YoutubeStopViewWS(userId, y)
}

func ViewNow(c *iris.Context) {
	userId := c.Get("token").(users.AccessToken).User.ID
	var yd YoutubeDonate
	c.ReadJSON(&yd)
	y, _ := json.Marshal(yd)
	YoutubeViewNowWS(userId, y)
}
