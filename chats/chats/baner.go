package chats

import (
	"errors"
	"os"

	"log"

	"time"

	"github.com/FireGM/chats/interfaces"
	"github.com/FireGM/chats/youtube"
)

func ban(chatname, channelID, userID, token string, t int) error {
	log.Println(chatname, channelID, userID, token, t)
	var chat interfaces.Bot
	clearChan := make(chan interfaces.Message, 1)
	ff := func(m interfaces.Message, b interfaces.Bot) {
		clearChan <- m
	}
	switch chatname {
	case "youtube":
	case "gplus":
		chat = youtube.NewWithAuth(ff, os.Getenv("YOUTUBE_KEY"), token)
	case "peka2tv":
		chat = connectPeka(ff, token)
	case "goodgame":
		chat = connectGoodGameToken(ff, token)
	case "twitch":
		chat = connectTwitch(ff, channelID, token, os.Getenv("TWITCH_KEY"))
	default:
		return errors.New("chat " + chatname + " not supported")
	}
	err := chat.Join(channelID)
	if err != nil {
		return err
	}
	err = chat.Timeout(channelID, userID, t)
	if err != nil {
		return err
	}
	timer := time.NewTimer(time.Second * 10)
clear:
	for {
		select {
		case msg := <-clearChan:
			if msg.GetChannelName() == channelID && msg.GetUID() == userID {
				break clear
			}
		case <-timer.C:
			timer.Stop()
			return errors.New("cant ban user " + userID)
		}
	}
	err = chat.Disconnect()
	if err != nil {
		return err
	}
	return nil
}

func sendMessageToChats(chats []StandardChatInfo, message string) []string {
	errorsR := make([]string, 0)
	for _, b := range chats {
		err := messageToChat(b.Chatname, b.ChannelID, b.Token, message)
		if err != nil {
			log.Println(err, b.Chatname)
			errorsR = append(errorsR, err.Error())
		}
	}
	return errorsR
}

func messageToChat(chatname, channelID, token, message string) error {
	log.Println(chatname, channelID, token, message)
	var chat interfaces.Bot
	clearChan := make(chan interfaces.Message, 1)
	ff := func(m interfaces.Message, b interfaces.Bot) {
		clearChan <- m
	}
	switch chatname {
	case "youtube":
	case "gplus":
		chat = youtube.NewWithAuth(ff, os.Getenv("YOUTUBE_KEY"), token)
	case "peka2tv":
		chat = connectPeka(ff, token)
	case "goodgame":
		chat = connectGoodGameToken(ff, token)
	case "twitch":
		chat = connectTwitch(ff, channelID, token, os.Getenv("TWITCH_KEY"))
	default:
		return errors.New("chat " + chatname + " not supported")
	}
	err := chat.Join(channelID)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	err = chat.SendMessageToChan(channelID, message)
	if err != nil {
		return err
	}
	timer := time.NewTimer(time.Second * 3)
clear:
	for {
		select {
		case msg := <-clearChan:
			log.Println(msg)
			if msg.GetChannelName() == channelID && msg.GetTextMessage() == message {
				break clear
			}
		case <-timer.C:
			timer.Stop()
			return errors.New("cant send message to channel " + channelID)
		}
	}
	err = chat.Disconnect()
	if err != nil {
		return err
	}
	return nil
}
