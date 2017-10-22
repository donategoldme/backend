package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"donategold.me/centrifugo"
	"donategold.me/db"
	"github.com/jinzhu/gorm"
)

const channelName = "youtubeWidget"

func init() {
	centrifugo.ChannelAllow(channelName)
	db.DB.AutoMigrate(&Youtube{}, &YoutubeDonate{})
}

type Youtube struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	Active      bool   `json:"active"`

	Likes    int `json:"likes"`
	Views    int `json:"views"`
	Duration int `json:"duration"`

	DeletedAt *time.Time `json:"-"`

	//Preload
	Videos []YoutubeDonate `json:"videos"`
}

func (y *Youtube) Create() {
	db.DB.Create(y)
	y.Videos = make([]YoutubeDonate, 0)
}

func (y *Youtube) Save() {
	var yt Youtube
	if db.DB.Where("id = ? and user_id = ?", y.ID, y.UserID).First(&yt); yt.ID == 0 {
		return
	}
	db.DB.Save(y)
	b, _ := json.Marshal(y)
	YoutubeSaveWS(y.UserID, b)
}

func (y *Youtube) AddVideo(data map[string]string) error {
	//d := make(map[string]string, 0)
	var yd YoutubeDonate
	yd.YoutubeID = y.ID
	yd.VideoId = getVideoIdFromUrl(data["url"])
	if yd.VideoId == "" {
		fmt.Println(3)
		return errors.New("No video url")
	}
	vd, err := getVideoDetails(yd.VideoId)
	if err != nil {
		fmt.Println(1)
		return err
	}
	err = ValidVideo(*y, vd, &yd)
	if err != nil {
		fmt.Println(2, err)
		return err
	}
	yd.Title = vd.Snippet.Title
	yd.UserID = y.UserID
	if _, ok := data["byUser"]; ok {
		yd.FromOwner = true
	} else {
		if nickname, ok := data["nickname"]; ok {
			yd.Nickname = nickname
		}
	}
	yd.Create()

	return nil
}

func (y *Youtube) ViewedVideo(yd YoutubeDonate) error {
	if yd.YoutubeID != y.ID {
		return errors.New("Not video")
	}
	yd.ViewedDone()
	b, err := json.Marshal(yd)
	if err != nil {
		return err
	}
	YoutubeViewdVideoWS(y.UserID, b)
	return nil
}

func GetYoutubeWidgetByUserAndId(userId uint, id uint) (Youtube, error) {
	var yw Youtube
	err := db.DB.Where("user_id = ? and id = ?", userId, id).First(&yw).Error
	if err != nil {
		return yw, err
	}
	return yw, nil
}

func GetActiveYWByUserID(userID uint, countVideo bool) (Youtube, bool) {
	var yw Youtube
	s := db.DB
	if countVideo {
		s = s.Preload("Videos", func(db *gorm.DB) *gorm.DB {
			db = db.Select("id, nickname, title, duration, likes, video_id, from_owner")
			db = db.Where("viewed = ? and (trans_id <> ? or from_owner = true)", false, 0)
			return db
		})
	}
	s.Where("user_id = ? and active = true", userID).First(&yw)
	if yw.ID == 0 {
		return yw, false
	}
	return yw, true
}

func GetYoutubeWidjetsByUser(user_id uint) []Youtube {
	var yw []Youtube
	db.DB.Preload("Videos", "viewed = ?", false).Where("user_id = ?", user_id).Order("-id").Find(&yw)
	return yw
}

func deleteWidgetByUserAndId(user_id uint, id int) error {
	err := db.DB.Where("id = ? and user_id = ?", id, user_id).Delete(Youtube{}).Error
	return err
}

func activing(userId uint, id int) error {
	tx := db.DB.Begin()
	if err := tx.Model(&Youtube{}).Where("user_id = ?", userId).Update("active", false).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&Youtube{}).Where("user_id = ? and id = ?", userId, id).Update("active", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

type YoutubeDonate struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	Viewed    bool       `json:"viewed"`
	Nickname  string     `json:"nickname"`
	Title     string     `json:"title"`
	Duration  int        `json:"duration"`
	Likes     int        `json:"likes"`
	Views     int        `json:"views"`
	VideoId   string     `json:"video_id"`
	YoutubeID uint       `json:"youtube_id"`
	TransID   uint       `json:"trans_id"`
	Payed     bool       `json:"-"`
	FromOwner bool       `json:"from_owner"`
	UserID    uint       `json:"user_id"`
	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (yd *YoutubeDonate) Create() error {
	db.DB.Create(yd)
	if yd.FromOwner {
		ydString, err := json.Marshal(&yd)
		if err != nil {
			return err
		}
		AddYoutubeWS(yd.UserID, ydString)
	}
	return nil
}

func (yd *YoutubeDonate) Save() {
	db.DB.Save(yd)
	if yd.TransID != 0 && yd.Viewed != true && yd.Payed {
		ydString, err := json.Marshal(&yd)
		if err != nil {
			return
		}
		AddYoutubeWS(yd.UserID, ydString)
	}
}

func (yd *YoutubeDonate) ViewedDone() {
	yd.Viewed = true
	db.DB.Model(yd).Update("viewed", true)
}

func GetYoutubeVideoByUserAndID(id uint) (YoutubeDonate, error) {
	var yd YoutubeDonate
	db.DB.Where("id = ?", id).First(&yd)
	if yd.ID == 0 {
		return yd, errors.New("No video")
	}
	return yd, nil
}

func CreateFromDonate(id uint, userID uint, transID uint, money int) error {
	// var yd YoutubeDonate
	yd, err := GetYoutubeVideoByUserAndID(id)
	if err != nil {
		return err
	}
	yt, err := GetYoutubeWidgetByUserAndId(userID, yd.YoutubeID)
	if err != nil {
		return err
	}
	if yt.Cost != money {
		return errors.New("Cost does not match")
	}
	// json.Unmarshal([]byte(data), &yd)
	yd.TransID = transID
	yd.Payed = true
	yd.Save()
	ydString, _ := json.Marshal(&yd)
	AddYoutubeWS(userID, ydString)
	return nil
}

func ValidVideo(y Youtube, vd VideoDetail, yd *YoutubeDonate) error {
	if l, err := strconv.Atoi(vd.Statistics.LikeCount); l < y.Likes || err != nil {
		log.Println("Need more likes")
		return errors.New("Need more likes")
	} else {
		yd.Likes = l
	}
	if v, err := strconv.Atoi(vd.Statistics.ViewCount); v < y.Views || err != nil {
		log.Println("Need more views")
		return errors.New("Need more views")
	} else {
		yd.Views = v
	}
	if d, err := getDuratonToSeconds(vd.ContentDetails.Duration); d > y.Duration || err != nil {
		log.Println("To long")
		return errors.New("To long")
	} else {
		yd.Duration = d
	}
	if vd.Snippet.LiveBroadcastContent == "live" {
		log.Println("It's live stream")
		return errors.New("It's live stream")
	}
	return nil
}
