package auth

import (
	"encoding/json"

	"fmt"

	"donategold.me/centrifugo"
	"donategold.me/components/users"
)

const channelName = "auth"

func init() {
	centrifugo.ChannelAllow(channelName)
}

func publishAddProvider(provider users.Provider) error {
	data, err := json.Marshal(&provider)
	if err != nil {
		return err
	}
	data = []byte(fmt.Sprintf(`{"type": "add_provider_auth", "provider": %s}`, data))
	return centrifugo.PublishChanAndMain(provider.UserID, channelName, data)
}
