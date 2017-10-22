package main

import (
	"log"

	"donategold.me/centrifugo"
	"donategold.me/server"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	centrifugo.Connect()
	server := server.Get()
	server.Listen(":80")

}
