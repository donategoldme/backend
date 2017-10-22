package main

import (
	_ "donategold.me/chats/centrifugo"
	"donategold.me/chats/chats"
	_ "donategold.me/chats/db"

	"log"

	"donategold.me/chats/polls"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/middleware/logger"
)

func main() {
	chats.ConnectChats()
	log.SetFlags(log.Llongfile | log.LstdFlags)

	server := iris.New()
	server.Adapt(iris.DevLogger())
	server.Adapt(httprouter.New())
	server.Use(logger.New())
	chats.Concat(server.Party("chats"))
	polls.Concat(server.Party("polls"))
	server.Listen(":80")
}
