package youtube

import (
	"fmt"

	"donategold.me/centrifugo"
)

func YoutubeSaveWS(userID uint, y []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "youtube_save", "youtube": %s}`, y))
	return centrifugo.PublishChanAndMain(userID, channelName, data)
}

func AddYoutubeWS(userId uint, yd []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "add_youtube_video", "video": %s}`,
		yd))
	return centrifugo.PublishChanAndMain(userId, channelName, data)
}

func YoutubeVideoPlayNowWS(userId uint, yd []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "youtube_video_play_now", "video": %s}`,
		yd))
	return centrifugo.PublishChanAndMain(userId, channelName, data)
}

func YoutubeViewdVideoWS(userId uint, yd []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "youtube_video_viewed", "video": %s}`,
		yd))
	return centrifugo.PublishChanAndMain(userId, channelName, data)
}

func YoutubeStopViewWS(userId uint, yd []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "youtube_stop_view", "video": %s}`,
		yd))
	return centrifugo.PublishChanAndMain(userId, channelName, data)
}

func YoutubeViewNowWS(userId uint, yd []byte) error {
	data := []byte(fmt.Sprintf(`{"type": "youtube_view_now", "video": %s}`,
		yd))
	return centrifugo.PublishChanAndMain(userId, channelName, data)
}
