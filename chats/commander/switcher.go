package commander

import (
	"log"
	"strings"

	"donategold.me/chats/polls"
	"github.com/centrifugal/gocent"
)

func CommandExist(s string) bool {
	return strings.HasPrefix(s, "!")
}

func CommandSwitcher(userID uint, text, username string, cgo *gocent.Client) {
	splitted := strings.SplitN(text, " ", 2)
	command, args := splitted[0], ""
	if len(splitted) == 2 {
		args = splitted[1]
	}
	log.Println(command)
	switch command {
	case "!vote":
		log.Println("!vote switcher")
		err := polls.VotePoll(userID, args, username, cgo)
		if err != nil {
			log.Println(err)
		}
	}
}
