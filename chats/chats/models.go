package chats

import (
	"errors"
	"html/template"
	"time"
)

type Message struct {
	UID           string        `json:"uid"`
	DisplayName   string        `json:"display_name"`
	ChannelID     string        `json:"channel_id"`
	ChatName      string        `json:"chat_name"`
	Moderator     bool          `json:"moderator"`
	ModeratorUrl  string        `json:"moderator_url,omitempty"`
	Subscriber    bool          `json:"subscriber"`
	SubscriberUrl string        `json:"subscriber_url,omitempty"`
	SmilesRender  template.HTML `json:"smiles_render"`
	FullRender    template.HTML `json:"full_render"`
}

type Chat struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	ChannelID string `json:"channel_id"`
	UserID    uint   `json:"user_id"`
	Slug      string `json:"slug"`

	DeletedAt *time.Time `json:"-"`
}

// struct for ban and send message
type StandardChatInfo struct {
	Chatname  string `json:"chat_name"`
	ChannelID string `json:"channel_id"`
	Token     string `json:"token"`
}

func reqErr(field string) error {
	return errors.New(field + " required")
}

func (b StandardChatInfo) Valid() error {
	if b.Chatname == "" {
		return reqErr("chat name")
	}
	if b.ChannelID == "" {
		return reqErr("channel id")
	}
	if b.Token == "" {
		return reqErr("token")
	}
	return nil
}

type BanUser struct {
	StandardChatInfo
	Time   int    `json:"time"`
	UserID string `json:"user_id"`
}

func (b *BanUser) Valid() error {
	if b.UserID == "" {
		return reqErr("user id")
	}
	if b.Time < 1 {
		b.Time = 1
	}
	return b.StandardChatInfo.Valid()
}

type MessageSend struct {
	Chats   []StandardChatInfo `json:"chats"`
	Message string             `json:"message"`
}

func (b *MessageSend) Valid() error {
	if b.Message == "" {
		return reqErr("message to send")
	}
	for _, chat := range b.Chats {
		if err := chat.Valid(); err != nil {
			return err
		}
	}
	return nil
}
