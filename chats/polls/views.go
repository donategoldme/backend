package polls

import (
	"encoding/json"
	"fmt"
	"log"

	"donategold.me/chats/centrifugo"

	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Post("/", addPoll)
	api.Get("/:id", getPolls)
	api.Delete("/:id", removePoll)
}

func getPolls(c *iris.Context) {
	userID, err := c.ParamInt("id")
	if err != nil {
		c.JSON(400, err)
		return
	}
	poll, exist := poolPolls.Get(uint(userID))
	if !exist {
		log.Printf("No poll for user %d\n", userID)
		return
	}
	data, err := json.Marshal(&poll)
	if err != nil {
		log.Println(err)
		return
	}
	data = []byte(fmt.Sprintf(`{"type": "get_polls_success", "poll": %s}`, data))
	centrifugo.PublishWS(uint(userID), "chats", data)
	c.JSON(200, poll)
}

func addPoll(c *iris.Context) {
	var poll PollFromJson
	c.ReadJSON(&poll)
	pollNew := poolPolls.NewPoll(poll)
	data, err := json.Marshal(&pollNew)
	if err != nil {
		return
	}
	data = []byte(fmt.Sprintf(`{"type": "save_polls_success", "poll": %s}`, data))
	centrifugo.PublishWS(poll.UserID, "chats", data)
	c.JSON(200, poll)
}

func removePoll(c *iris.Context) {
	userID, err := c.ParamInt("id")
	if err != nil {
		return
	}
	poolPolls.Clear(uint(userID))
	c.WriteString("ok")
}
