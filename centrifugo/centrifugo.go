package centrifugo

import (
	"os"
	"time"

	"strings"

	"sync"

	"github.com/centrifugal/gocent"
)

var Cgo *gocent.Client

func Connect() {
	Cgo = gocent.NewClient("http://"+os.Getenv("CENTRIFUGO_SERVICE")+":8000",
		os.Getenv("CENTRIFUGO_SECRET"), 5*time.Second)
}

type channelsAllow struct {
	channels map[string]bool
	sync.RWMutex
}

func (c *channelsAllow) Add(ch string) {
	c.Lock()
	defer c.Unlock()
	c.channels[ch] = true
}

func (c *channelsAllow) Check(ch string) bool {
	c.RLock()
	defer c.RUnlock()
	return c.channels[ch]
}

var Channels = channelsAllow{map[string]bool{"": true}, sync.RWMutex{}}

//ChannelsAllow add new channel for subscribe
func ChannelAllow(ch string) {
	Channels.Add(ch)
}

//CheckAllowChannel function exiting channel in channels
//for subscribe
func CheckAllowChannel(ch string) bool {
	channel := strings.Split(ch, "/")
	if len(channel) != 2 {
		return false
	}
	return Channels.Check(channel[1])
}
