package standard

import (
	"fmt"

	"donategold.me/components/users"
	"donategold.me/speecher"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Get("/", GetPreferences)
	api.Post("/", Create)
	api.Get("/donates", GetDonates)
	api.Post("/donates", AddDonate)
	api.Post("/donates/:id", ViewedDonate)
	api.Get("/donates/:id", TextToSpeech)
	api.Post("/widget/:id", Save)
	api.Delete("/widget/:id", Delete)
	api.Get("/paypage", getPaypages)
	api.Post("/paypage", createPaypage)
	api.Post("/paypage/:id", savePaypage)
	api.Delete("/paypage/:id", deletePaypage)
}

func Create(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	s := CreateWithDefault(userID)
	c.JSON(200, s)
}

func GetPreferences(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	s := GetPreferenceByUserId(userID)
	c.JSON(200, s)
}

func GetDonates(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	donates := GetDonatesByUser(userID)
	c.JSON(200, donates)
}

func AddDonate(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	var d StandardDonate
	c.ReadJSON(&d)
	d.UserID = userID
	d.FromOwner = true
	d.Create()
	c.JSON(200, d)
}

func ViewedDonate(c *iris.Context) {
	id, err := c.ParamInt("id")
	if err != nil || id == 0 {
		c.EmitError(400)
		return
	}
	userID := c.Get("token").(users.AccessToken).User.ID
	s, exist := GetDonateByUserAndId(userID, uint(id))
	if !exist {
		c.EmitError(404)
		return
	}
	if s.Viewed {
		return
	}
	s.Viewed = true
	s.Save()
	c.JSON(200, s)
}

func Save(c *iris.Context) {
	id, err := c.ParamInt("id")
	if err != nil || id == 0 {
		c.EmitError(400)
		return
	}
	userID := c.Get("token").(users.AccessToken).User.ID
	s, exist := GetPreferenceByUserIdAndId(userID, id)
	if !exist {
		c.EmitError(404)
		return
	}
	c.ReadJSON(&s)
	s.UserID = userID
	s.ID = uint(id)
	s.Save()
	c.JSON(200, s)
}

func Delete(c *iris.Context) {
	id, err := c.ParamInt("id")
	if err != nil {
		c.EmitError(400)
		return
	}
	userID := c.Get("token").(users.AccessToken).User.ID
	s, exist := GetPreferenceByUserIdAndId(userID, id)
	if !exist {
		c.EmitError(404)
		return
	}
	s.Delete()
	c.JSON(200, s.ID)
}

func TextToSpeech(c *iris.Context) {
	id, err := c.ParamInt("id")
	if err != nil {
		c.EmitError(400)
		return
	}
	userID := c.Get("token").(users.AccessToken).User.ID
	_, exist := GetDonateByUserAndId(userID, uint(id))
	if !exist {
		c.EmitError(404)
		return
	}
	text := c.URLParam("text")
	err = speecher.Client.SaveToAudio(text, fmt.Sprintf(voicePathFormat, id), 0777)
	if err != nil {
		c.EmitError(400)
		return
	}
	c.JSON(200, "ok")
}

func getPaypages(c *iris.Context) {
	user := c.Get("token").(users.AccessToken)
	paypages := GetByUserID(user.User.ID)
	c.JSON(200, paypages)
}

func createPaypage(c *iris.Context) {
	var p Paypage
	user := c.Get("token").(users.AccessToken)
	c.ReadJSON(&p)
	p.UserID = user.User.ID
	p.ID = 0
	p.Create()
	c.JSON(201, p)
}

func savePaypage(c *iris.Context) {
	var p Paypage
	user := c.Get("token").(users.AccessToken)
	c.ReadJSON(&p)
	id, err := c.ParamInt("id")
	if err != nil {
		c.JSON(404, "Need integer id")
		return
	}
	p.UserID = user.User.ID
	p.ID = uint(id)
	_, exist := GetByUserIDAndID(p.UserID, p.ID)
	if !exist {
		c.JSON(404, "not found")
		return
	}
	p.Save()
	c.JSON(200, p)
}

func deletePaypage(c *iris.Context) {
	user := c.Get("token").(users.AccessToken)
	id, err := c.ParamInt("id")
	if err != nil {
		c.JSON(404, "Not found")
		return
	}
	DeleteByUserIDAndID(user.User.ID, uint(id))
	c.JSON(200, "Deleted")
}
