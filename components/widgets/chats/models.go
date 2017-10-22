package chats

import (
	"errors"
	"time"

	"donategold.me/db"
)

func init() {
	db.DB.AutoMigrate(&Chat{}, &Pref{})
}

//border
//color message/nickname
//color background
//font size and type
//padding bot
//count of messages
//badges
type Pref struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	ColorMessage string `json:"color_message"`
	ColorNicks   string `json:"color_nicks"`
	ColorBG      string `json:"color_bg"`
	FontSize     int    `json:"font_size"`
	FontFamily   string `json:"font_family"`
	Padding      int    `json:"padding"`
	MarginBot    int    `json:"margin_bot"`
	BorderRadius int    `json:"border_radius"`
	Badges       bool   `json:"badges"`
	CSS          string `json:"css"`

	DeletedAt *time.Time `json:"-"`
}

func (c *Pref) Create() {
	db.DB.Create(c)
}

func (c *Pref) Save() {
	db.DB.Save(c)
	savePrefWS(c)
}

func getPrefByUser(userID uint) Pref {
	var pref Pref
	db.DB.Where("user_id = ?", userID).First(&pref)
	if pref.ID == 0 {
		pref.UserID = userID
		db.DB.Create(&pref)
	}
	return pref
}

type ChatUrl struct {
	Url string `json:"url"`
}

type Chat struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	ChannelID string `json:"channel_id"`
	UserID    uint   `json:"user_id"`
	Slug      string `json:"slug"`

	DeletedAt *time.Time `json:"-"`
}

func (c *Chat) Create() error {
	var chat Chat
	db.DB.Where("user_id = ? and type = ? and channel_id = ?", c.UserID, c.Type, c.ChannelID).First(&chat)
	if chat.ID != 0 {
		return errors.New("already exist this chat")
	}
	if err := db.DB.Create(c).Error; err != nil {
		return err
	}
	return subscribeChatWS(*c)
}

func (c *Chat) Save() {
	db.DB.Save(c)
}

func (c *Chat) Delete() {
	db.DB.Delete(c)
}

func GetChatsByUser(userID uint) []Chat {
	var chats []Chat
	db.DB.Where("user_id = ?", userID).Find(&chats)
	return chats
}

func GetChatByIdAndUserId(id, userID uint) Chat {
	var chat Chat
	db.DB.Where("id = ? and user_id = ?", id, userID).First(&chat)
	return chat
}

func RemoveSubChat(id int) {
	db.DB.Delete(&Chat{}, "id = ?", id)
}

type standardInfo struct {
	Chatname    string `json:"chat_name"`
	ChannelID   string `json:"channel_id"`
	ProviderUID string `json:"provider_uid"`
	Token       string `json:"token"`
}

type BanUser struct {
	standardInfo
	Time   int    `json:"time"`
	UserID string `json:"user_id"`
}

type SendMessage struct {
	Chats   []standardInfo `json:"chats"`
	Message string         `json:"message"`
}

type Poll struct {
	UserID   uint     `json:"user_id"`
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}
