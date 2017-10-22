package chats

import (
	"log"

	"donategold.me/chats/chatters"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Post("/add", addChat)
	api.Post("/remove", removeChat)
	api.Post("/ban", banUser)
	api.Post("/send", sendMessage)
}

func addChat(c *iris.Context) {
	var chat Chat
	c.ReadJSON(&chat)
	log.Println(chat)
	if v, ok := chatsMap[chat.Type]; ok {
		log.Println(ok, chat.ChannelID)
		err := v.Join(chat.ChannelID)
		if err != nil {
			log.Println(err)
			c.JSON(500, err.Error())
			return
		}
		chatters.Add(chat.Type+"/"+chat.ChannelID, chat.UserID)
		c.JSON(200, "ok")
		return
	}
	log.Println("chat not supported")
	c.JSON(404, "chat not supported")
}

func removeChat(c *iris.Context) {
	var chat Chat
	c.ReadJSON(&chat)
	subs := chatters.Remove(chat.Type+"/"+chat.ChannelID, chat.UserID)
	if subs == 0 {
		chatsMap["youtube"].Leave(chat.ChannelID)
	}
	c.JSON(200, "ok")
}

func banUser(c *iris.Context) {
	var b BanUser
	c.ReadJSON(&b)
	if err := b.Valid(); err != nil {
		c.JSON(400, err.Error())
		return
	}
	if err := ban(b.Chatname, b.ChannelID, b.UserID, b.Token, b.Time); err != nil {
		c.JSON(400, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func sendMessage(c *iris.Context) {
	var b MessageSend
	c.ReadJSON(&b)
	if err := b.Valid(); err != nil {
		c.JSON(400, err.Error())
		return
	}
	if err := sendMessageToChats(b.Chats, b.Message); len(err) > 0 {
		c.JSON(400, err)
		return
	}
	c.JSON(200, "ok")
}
