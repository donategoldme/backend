package chats

import (
	"encoding/json"
	"fmt"
	"os"

	"donategold.me/chats/centrifugo"
	"donategold.me/chats/chatters"
	"donategold.me/chats/commander"
	"donategold.me/chats/db"
	"donategold.me/chats/polls"

	"log"

	"github.com/FireGM/chats"
	"github.com/FireGM/chats/goodgame"
	"github.com/FireGM/chats/interfaces"
	"github.com/FireGM/chats/peka2tv"
	"github.com/FireGM/chats/twitch"
	"github.com/FireGM/chats/youtube"
	"github.com/centrifugal/gocent"
)

var chatsMap = map[string]interfaces.Bot{}

func ConnectChats() {
	ch := make(chan interfaces.Message, 1000)
	poolPolls := polls.GetPoolOfPolls()
	handler := chats.MakerHandlers(ch)
	chatsMap["peka2tv"] = connectPeka(handler, os.Getenv("PEKA2TV_TOKEN"))
	chatsMap["goodgame"] = connectGoodGame(handler, os.Getenv("GOODGAME_LOGIN"), os.Getenv("GOODGAME_PASS"))
	chatsMap["twitch"] = connectTwitch(handler, os.Getenv("TWITCH_NAME"),
		os.Getenv("TWITCH_OAUTH"), os.Getenv("TWITCH_KEY"))
	chatsMap["youtube"] = connectYoutube(handler, os.Getenv("YOUTUBE_KEY"))
	loadChatsJoins()
	go publisher(ch, centrifugo.Cgo, poolPolls)
}

func publisher(ch chan interfaces.Message, cgo *gocent.Client, pp *polls.PoolOfPolls) {
	for m := range ch {
		message := getMessage(m)
		data := messageToData(message, m.IsClearMessage())
		for _, userID := range chatters.GetSubsOfChan(m.GetChatName() + "/" + m.GetChannelName()) {
			_, err := centrifugo.Cgo.Publish(fmt.Sprintf("$%d/chats", userID), data)
			if err != nil {
				log.Println(err)
			}
			if commander.CommandExist(m.GetTextMessage()) {
				commander.CommandSwitcher(userID, m.GetTextMessage(), m.GetUserFrom(), cgo)
			}
		}
	}
}

func getMessage(m interfaces.Message) Message {
	var message Message
	message.ChannelID = m.GetChannelName()
	message.DisplayName = m.GetUserFrom()
	message.SmilesRender = m.GetRenderSmiles()
	message.FullRender = m.GetRenderFullHTML()
	message.UID = m.GetUID()
	message.ChatName = m.GetChatName()
	message.Moderator, message.ModeratorUrl = m.IsModerator()
	message.Subscriber, message.SubscriberUrl = m.IsSubscriber()
	return message
}

func messageToData(m Message, clear bool) []byte {
	str, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	typeOfMessage := ""
	if clear {
		typeOfMessage = "clear_message"
	} else {
		typeOfMessage = "new_message"
	}
	return []byte(fmt.Sprintf(`{"type": "%s", "message": %s}`, typeOfMessage, str))
}

func loadChatsJoins() {
	var chatsLoad []Chat
	db.DB.Find(&chatsLoad)
	for _, c := range chatsLoad {
		err := chatsMap[c.Type].Join(c.ChannelID)
		if err != nil {
			log.Println(err)
		}
		chatters.Add(c.Type+"/"+c.ChannelID, c.UserID)
	}
}

func connectPeka(handle func(interfaces.Message, interfaces.Bot), token string) *peka2tv.Bot {
	chat := peka2tv.New(handle)
	chat.Connect()
	chat.LoginByToken(token)
	return chat
}

func connectGoodGame(handle func(interfaces.Message, interfaces.Bot), login, pass string) *goodgame.Bot {
	chat := goodgame.New(handle)
	chat.Connect()
	chat.LoginByPass(login, pass)
	return chat
}

func connectGoodGameToken(handle func(interfaces.Message, interfaces.Bot), token string) *goodgame.Bot {
	chat := goodgame.New(handle)
	chat.Connect()
	chat.LoginByToken(token)
	return chat
}

func connectTwitch(handle func(interfaces.Message, interfaces.Bot), name, oauth, key string) *twitch.Bot {
	chat := twitch.NewWithRender(name, oauth, key, handle)
	err := chat.Connect()
	if err != nil {
		log.Println(err)
	}
	return chat
}

func connectYoutube(handle func(interfaces.Message, interfaces.Bot), key string) *youtube.Bot {
	chat := youtube.New(handle, key)
	return chat
}
