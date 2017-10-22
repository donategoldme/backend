package chatters

import (
	"sync"
)

func New() *Chatter {
	chatter := Chatter{}
	chatter.init()
	return &chatter
}

type Chatter struct {
	channels map[string][]uint
	sync.RWMutex
}

func (c *Chatter) init() {
	c.channels = make(map[string][]uint)
}

func (c *Chatter) Add(channel string, user uint) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.channels[channel]; !ok {
		c.channels[channel] = []uint{}
	}
	for _, id := range c.channels[channel] {
		if id == user {
			return
		}
	}
	c.channels[channel] = append(c.channels[channel], user)
}

func (c *Chatter) Remove(channel string, user uint) int {
	c.Lock()
	defer c.Unlock()
	if v, ok := c.channels[channel]; ok {
		for i, userID := range v {
			if userID == user {
				c.channels[channel] = append(v[:i], v[i+1:]...)
			}
		}
	}
	return len(c.channels[channel])
}

func (c *Chatter) SubsUser(userID uint) []string {
	c.RLock()
	defer c.RUnlock()
	subs := []string{}
	for k, v := range c.channels {
		for _, id := range v {
			if id == userID {
				subs = append(subs, k)
			}
		}
	}
	return subs
}

func (c *Chatter) GetSubsOfChan(channel string) []uint {
	c.RLock()
	defer c.RUnlock()
	return c.channels[channel]
}
