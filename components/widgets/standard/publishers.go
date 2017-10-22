package standard

import (
	"fmt"

	"donategold.me/centrifugo"
)

const channelName string = "standard"

func init() {
	centrifugo.ChannelAllow(channelName)
}

func sendAddStandardDonateWS(userID uint, sd []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "add_standard_donate", "donate": %s}`, sd))
	return centrifugo.PublishChanAndMain(userID, channelName, data)
}

func sendSaveStandardDonateWS(userID uint, sd []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "save_standard_donate", "donate": %s}`, sd))
	return centrifugo.PublishChanAndMain(userID, channelName, data)
}
func sendSaveStandardPrefsWS(userID uint, standard []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "standard_prefs_save", "standard": %s}`, standard))
	return centrifugo.PublishChanAndMain(userID, channelName, data)
}

func sendCreatePaypageWS(userID uint, paypage []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "paypage_create", "paypage": %s}`, paypage))
	return centrifugo.PublishChanAndMain(userID, channelName, data)
}

func sendSavePaypageWS(userID uint, paypage []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "paypage_save", "paypage": %s}`, paypage))
	return centrifugo.PublishChanAndMain(userID, channelName, data)
}

func sendDeletePaypageWS(userID uint, paypage []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "paypage_delete", "paypage": %s}`, paypage))
	return centrifugo.PublishChanAndMain(userID, channelName, data)
}
