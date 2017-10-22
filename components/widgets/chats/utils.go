package chats

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"net/http"
	"time"

	"bytes"

	"os"

	"io/ioutil"

	"github.com/FireGM/chats/goodgame"
	"github.com/FireGM/chats/peka2tv"
)

const channelInfoYoutubeURL = `https://www.googleapis.com/youtube/v3/channels`

const subURL = "http://chats/chats/add"
const unSubURL = "http://chats/chats/remove"
const banURL = "http://chats/chats/ban"
const sendMessageURL = "http://chats/chats/send"
const getPollURL = "http://chats/polls"
const addPollURL = "http://chats/polls/"
const removePollURL = "http://chats/polls/"

var client = http.Client{Timeout: time.Second * 30}

func getChat(urlStr string, userID uint) (Chat, error) {
	var chat Chat
	if urlStr == "" {
		return chat, errors.New("No url")
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return chat, err
	}
	// fmt.Println(u.Hostname())
	switch u.Hostname() {
	case "www.twitch.tv", "twitch.tv":
		p := strings.Split(u.EscapedPath(), "/")
		if len(p) < 2 {
			return chat, errors.New("No channel name")
		}
		chat = Chat{Type: "twitch", ChannelID: p[1], UserID: userID, Slug: p[1]}
	case "www.goodgame.ru", "goodgame.ru":
		p := strings.Split(u.EscapedPath(), "/")
		if len(p) < 3 {
			return chat, errors.New("No channel name")
		}
		chanID, err := goodgame.GetStreamInfo(p[2])
		if err != nil {
			return chat, err
		}
		chat = Chat{Type: "goodgame", ChannelID: strconv.Itoa(chanID), UserID: userID, Slug: p[2]}
	case "www.peka2.tv", "peka2.tv":
		p := strings.Split(u.EscapedPath(), "/")
		if len(p) < 2 {
			return chat, errors.New("No channel name")
		}
		chanID, err := peka2tv.GetUserIdBySlug(p[1])
		if err != nil {
			return chat, err
		}
		chat = Chat{Type: "peka2tv", ChannelID: "stream/" + strconv.Itoa(chanID),
			UserID: userID, Slug: p[1]}
	case "www.youtube.com", "youtube.com":
		p := strings.Split(u.EscapedPath(), "/")
		if len(p) < 3 {
			return chat, errors.New("No channel id")
		}
		chanTitle, err := getChannelDisplayName(p[2])
		if err != nil {
			return chat, err
		}
		chat = Chat{Type: "youtube", ChannelID: p[2], UserID: userID, Slug: chanTitle}
	default:
		return chat, errors.New("unsupported site")
	}
	return chat, nil
}

func subscribeUserToChan(chat Chat) error {
	data, err := json.Marshal(&chat)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", subURL, bytes.NewBuffer(data))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("Repeat. Probably not found stream on channel")
	}
	return nil
}

func unsubscribeUserFromChan(chat Chat) error {
	data, err := json.Marshal(&chat)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", unSubURL, bytes.NewBuffer(data))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("Repeat")
	}
	return nil
}

func getPollsService(userID uint) error {
	// log.Println(fmt.Sprintf("%s/%d", getPollURL, userID))
	res, err := client.Get(fmt.Sprintf("%s/%d", getPollURL, userID))
	if err != nil {
		return err
	}
	// b, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }
	// log.Println(string(b))
	res.Body.Close()
	return nil
}

func addPollService(p Poll) error {
	data, err := json.Marshal(&p)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", addPollURL, bytes.NewBuffer(data))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("Repeat")
	}
	return nil
}

func removePollService(userID uint) error {
	req, _ := http.NewRequest("DELETE", removePollURL+strconv.Itoa(int(userID)), nil)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("Repeat")
	}
	return nil
}

type ChanInfoResp struct {
	Items []struct {
		Snippet struct {
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

func getChannelDisplayName(channelId string) (string, error) {
	req, err := http.NewRequest("GET", channelInfoYoutubeURL, nil)
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("id", channelId)
	values.Set("part", "snippet")
	values.Set("key", os.Getenv("YOUTUBE_KEY"))
	req.URL.RawQuery = values.Encode()
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var resp ChanInfoResp
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return "", err
	}
	if len(resp.Items) < 1 {
		return "", errors.New("Cant find channel")
	}
	return resp.Items[0].Snippet.Title, nil
}

func banUserOnChan(banUser BanUser) error {
	data, err := json.Marshal(&banUser)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", banURL, bytes.NewBuffer(data))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New("Cant read response from chatservice")
		}
		return errors.New(string(b))
	}
	return nil
}

func sendMessageOnChans(sender SendMessage) error {
	data, err := json.Marshal(&sender)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", sendMessageURL, bytes.NewBuffer(data))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New("Cant read response from chatservice")
		}
		return errors.New(string(b))
	}
	return nil
}
