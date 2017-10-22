package subscribers

import (
	"encoding/json"

	"donategold.me/centrifugo"
)

const channelName = "subscribers"

func init() {
	centrifugo.ChannelAllow(channelName)
}

func subscribeChatWS() error {
	data, err := json.Marshal(&chat)
	if err != nil {
		return err
	}

	return centrifugo.PublishChan(chat.UserID, channelName, data)
}
