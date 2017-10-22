package polls

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/centrifugal/gocent"
)

func VotePoll(userID uint, args string, username string, cgo *gocent.Client) error {
	spl := strings.SplitN(args, " ", 1)
	answerIndex, err := strconv.Atoi(spl[0])
	if err != nil {
		log.Println("!vote")
		return err
	}
	log.Println("!vote")
	p, err := poolPolls.Vote(userID, answerIndex, username)
	if err != nil {
		return err
	}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	data = []byte(fmt.Sprintf(`{"type": "get_polls_success", "poll": %s}`, data))
	_, err = cgo.Publish(fmt.Sprintf("$%d/chats", userID), data)
	return err
}
