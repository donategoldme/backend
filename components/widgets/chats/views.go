package chats

import (
	"log"
	"strconv"

	"donategold.me/components/auth"
	"donategold.me/components/users"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Get("/", getChats)

	api.Get("/pref", getPref)
	api.Post("/pref", savePref)

	api.Post("/conn/add", addChat)
	api.Post("/conn/remove", removeChat)

	api.Get("/polls", getPolls)
	api.Post("/polls", addPoll)
	api.Delete("/polls", removePoll)
	api.Post("/polls/view", viewPollOnScreen)

	api.Post("/ban", banUserView)
	api.Post("/send", sendMessageView)
}

func getChats(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	chats := GetChatsByUser(userID)
	c.JSON(200, chats)
}

func getPref(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	pref := getPrefByUser(userID)
	c.JSON(200, pref)
}

func savePref(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	pref := getPrefByUser(userID)
	var p Pref
	c.ReadJSON(&p)
	p.ID = pref.ID
	p.UserID = userID
	p.Save()
	c.JSON(200, p)
}

func addChat(c *iris.Context) {
	var chatURL ChatUrl
	c.ReadJSON(&chatURL)
	if chatURL.Url == "" {
		c.EmitError(404)
	}
	userID := c.Get("token").(users.AccessToken).User.ID
	chat, err := getChat(chatURL.Url, userID)
	if err != nil {
		log.Println(err)
		log.Println(chatURL)
		c.JSON(500, err.Error())
		return
	}
	err = chat.Create() //todo: get without create
	if err != nil {
		log.Println(err)
		c.JSON(500, err.Error())
		return
	}
	err = subscribeUserToChan(chat)
	if err != nil {
		log.Println(err)
		chat.Delete()
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, chat)
}

func removeChat(c *iris.Context) {
	var chat Chat
	c.ReadJSON(&chat)
	userID := c.Get("token").(users.AccessToken).User.ID
	if chat.UserID != userID {
		c.EmitError(404)
	}
	chat = GetChatByIdAndUserId(chat.ID, userID)
	err := unsubscribeUserFromChan(chat)
	if err != nil {
		log.Println(err)
		c.JSON(500, err.Error())
		return
	}
	err = unsubscribeUserWS(chat)
	if err != nil {
		log.Println(err)
		c.JSON(500, err.Error())
	}
	chat.Delete()
	c.JSON(200, chat)
}

func addPoll(c *iris.Context) {
	var p Poll
	err := c.ReadJSON(&p)
	if err != nil {
		c.JSON(400, err)
		return
	}
	userID := c.Get("token").(users.AccessToken).User.ID
	p.UserID = userID
	err = addPollService(p)
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.WriteString("ok")
}

func getPolls(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	err := getPollsService(userID)
	if err != nil {
		log.Println("Error: " + err.Error())
	}
}

func removePoll(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	err := removePollService(userID)
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.WriteString("ok")
}

func viewPollOnScreen(c *iris.Context) {
	userID := c.Get("token").(users.AccessToken).User.ID
	view, err := strconv.ParseBool(c.URLParam("view"))
	if err != nil {
		c.JSON(400, err)
		return
	}
	err = viewPollScreenWS(userID, view)
	if err != nil {
		c.JSON(400, err)
		return
	}
	if view {
		err = getPollsService(userID)
		if err != nil {
			log.Println("Error: " + err.Error())
			c.JSON(400, err)
			return
		}
	}
	c.WriteString("ok")
}

func banUserView(c *iris.Context) {
	var banUser BanUser
	err := c.ReadJSON(&banUser)
	if err != nil {
		c.JSON(500, "server error")
		return
	}
	user := c.Get("token").(users.AccessToken)
	provider, err := users.GetProviderByTypeAndUserIDAndUID(banUser.Chatname, user.User.ID, banUser.ProviderUID)
	if err != nil {
		log.Println(err)
		c.JSON(400, err.Error())
		return
	}
	_, err = auth.RefreshToken(&provider)
	if err != nil {
		log.Println(err)
		c.JSON(500, err.Error())
		return
	}
	banUser.Token = provider.AccessToken
	err = banUserOnChan(banUser)
	if err != nil {
		log.Println(err)
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func sendMessageView(c *iris.Context) {
	var sender SendMessage
	err := c.ReadJSON(&sender)
	if err != nil {
		c.JSON(500, "server error")
		return
	}
	log.Println(sender)
	user := c.Get("token").(users.AccessToken)
	for i := range sender.Chats {
		provider, err := users.GetProviderByTypeAndUserIDAndUID(sender.Chats[i].Chatname, user.User.ID, sender.Chats[i].ProviderUID)
		if err != nil {
			log.Println(err)
			c.JSON(400, err.Error())
			return
		}
		_, err = auth.RefreshToken(&provider)
		if err != nil {
			log.Println(err)
			c.JSON(500, err.Error())
			return
		}
		sender.Chats[i].Token = provider.AccessToken
	}

	err = sendMessageOnChans(sender)
	if err != nil {
		log.Println(err)
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}
