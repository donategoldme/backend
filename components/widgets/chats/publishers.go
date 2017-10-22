package chats

import (
	"encoding/json"
	"fmt"

	"donategold.me/centrifugo"
)

const channelName string = "chats"

func init() {
	centrifugo.ChannelAllow(channelName)
}

func subscribeChatWS(chat Chat) error {
	data, err := json.Marshal(&chat)
	if err != nil {
		return err
	}
	data = []byte(fmt.Sprintf(`{"type": "chat_subscribe", "chat": %s}`, data))
	return centrifugo.PublishChan(chat.UserID, channelName, data)
}

func unsubscribeUserWS(chat Chat) error {
	data, err := json.Marshal(&chat)
	if err != nil {
		return err
	}
	data = []byte(fmt.Sprintf(`{"type": "chat_unsubscribe", "chat": %s}`, data))
	return centrifugo.PublishChan(chat.UserID, channelName, data)
}

func savePrefWS(cp *Pref) error {
	data, err := json.Marshal(cp)
	if err != nil {
		return err
	}
	data = []byte(fmt.Sprintf(`{"type": "chats_pref_save", "chats_pref": %s}`, data))
	return centrifugo.PublishChanAndMain(cp.UserID, channelName, data)
}

func viewPollScreenWS(userID uint, view bool) error {
	data := []byte(fmt.Sprintf(`{"type":"polls_view_on_screen", "view": %t}`, view))
	return centrifugo.PublishChan(userID, channelName, data)
}
