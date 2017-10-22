package public

import (
	"fmt"
	"log"

	"donategold.me/components/users"
	"donategold.me/components/widgets/standard"
	"donategold.me/components/widgets/youtube"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Post("/youtube", createYoutube)
	api.Post("/standard", createStandard)
	api.Get("/youtube", getYoutubeInfo)
	api.Get("/standard", getStandardInfo)
}

// пресоздание видео для ютуба, но не отсылается пользователю до оплаты
func createYoutube(c *iris.Context) {
	var data Youtube
	c.ReadJSON(&data)
	user, exist := users.GetUserByUsername(data.Username)
	if !exist {
		c.JSON(404, "Username not found")
		return
	}
	youtube, exist := youtube.GetActiveYWByUserID(user.ID, false)
	if !exist {
		c.JSON(404, fmt.Sprintf("Active youtube widget for %s not found", user.Username))
		return
	}
	err := youtube.AddVideo(map[string]string{"url": data.VideoURL, "nickname": data.Nickname})
	if err != nil {
		log.Println(err.Error)
		c.JSON(500, "Server error")
		return
	}
	c.JSON(201, "ok")
}

// пресоздание стандартного доната
func createStandard(c *iris.Context) {
	var data Standard
	c.ReadJSON(&data)
	user, exist := users.GetUserByUsername(data.Username)
	if !exist {
		c.JSON(404, "Username not found")
		return
	}
	sd, exist := standard.GetActiveStandardByUserID(user.ID)
	if !exist {
		c.JSON(404, fmt.Sprintf("Active standard widget for %s not found", user.Username))
		return
	}
	if data.Money < sd.Cost {
		c.JSON(400, fmt.Sprintf("Not does match. Need more then %d", sd.Cost))
		return
	}
	d := standard.CreateStandardDonate(data.Nickname, user.ID, data.Message, data.Money, false)
	// log.Println(d)
	c.JSON(201, d)
}

func getYoutubeInfo(c *iris.Context) {
	username := c.URLParam("username")
	if username == "" {
		c.JSON(404, "Username required")
		return
	}
	user, exist := users.GetUserByUsername(username)
	if !exist {
		c.JSON(404, "Username not found")
		return
	}
	sd, exist := youtube.GetActiveYWByUserID(user.ID, true)
	if !exist {
		c.JSON(404, fmt.Sprintf("Active youtube widget for %s not found", user.Username))
		return
	}
	sd.ID = 0
	sd.UserID = 0
	c.JSON(200, sd)
}

func getStandardInfo(c *iris.Context) {
	username := c.URLParam("username")
	if username == "" {
		c.JSON(404, "Username required")
		return
	}
	user, exist := users.GetUserByUsername(username)
	if !exist {
		c.JSON(404, "Username not found")
		return
	}
	paypage, exist := standard.GetActivePaypageByUserID(user.ID)
	if !exist {
		c.JSON(404, fmt.Sprintf("Active standard paypage for %s not found", user.Username))
		return
	}
	cost, voiceCost, volume, exist := standard.GetParamsOfActiveWidgetByUserID(user.ID)
	if !exist {
		c.JSON(404, fmt.Sprintf("Active standard widget for %s not found", user.Username))
		return
	}
	paypage.ID = 0
	data := map[string]interface{}{"paypage": paypage, "cost": cost, "voice_cost": voiceCost,
		"volume": volume}
	c.JSON(200, data)
}
