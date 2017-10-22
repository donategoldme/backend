package centrifugo

import (
	"fmt"
	"os"
	"time"

	"github.com/centrifugal/gocent"
)

var Cgo *gocent.Client = Connect()

func Connect() *gocent.Client {
	client := gocent.NewClient("http://"+os.Getenv("CENTRIFUGO_SERVICE")+":8000",
		os.Getenv("CENTRIFUGO_SECRET"), 5*time.Second)
	return client
}

func PublishWS(userID uint, channel string, data []byte) error {
	_, err := Cgo.Publish(fmt.Sprintf("$%d/%s", userID, channel), data)
	return err
}
