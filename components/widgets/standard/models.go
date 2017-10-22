package standard

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"donategold.me/db"
)

const voicePath string = "/var/uploads/voice/"
const voicePathFormat string = voicePath + "%d.mp3"

func init() {
	os.MkdirAll(voicePath, 0777)
	db.DB.AutoMigrate(&StandardDonate{}, &Standard{})
}

type Standard struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	UserID       uint   `json:"user_id"`
	Cost         int    `json:"cost"` //==min_cost
	VoiceCost    int    `json:"voice_cost"`
	ViewTime     int    `json:"view_time"` //time on screen
	Pic          string `json:"pic"`
	Sound        string `json:"sound"`
	SoundVolume  int    `json:"sound_volume"`
	ViewTemplate string `json:"view_template"` // nickname amount message

	DeletedAt *time.Time `json:"-"`
}

func (s *Standard) Create() {
	db.DB.Create(s)
}

//Save to db. If pref active do another standards preference to deactive
//and pushed to centrifugo
func (s *Standard) Save() {
	if s.Active == true {
		tx := db.DB.Begin()
		tx.Model(&Standard{}).Where("user_id = ?", s.UserID).Update("active", false)
		tx.Save(s)
		tx.Commit()
	} else {
		db.DB.Save(s)
	}
	b, _ := json.Marshal(s)
	sendSaveStandardPrefsWS(s.UserID, b)
}

func (s *Standard) Delete() {
	db.DB.Delete(s)
}

func GetActiveStandardByUserID(userID uint) (Standard, bool) {
	var s Standard
	db.DB.Where("user_id = ? and active = ?", userID, true).First(&s)
	if s.ID == 0 {
		return s, false
	}
	return s, true
}

func GetParamsOfActiveWidgetByUserID(userID uint) (cost int, voice int, volume bool, exist bool) {
	s, exist := GetActiveStandardByUserID(userID)
	return s.Cost, s.VoiceCost, s.SoundVolume != 0, exist
}

//GetPreferenceByUserId get all standards preference for user
func GetPreferenceByUserId(id uint) []Standard {
	var s []Standard
	db.DB.Where("user_id = ?", id).Order("id").Find(&s)
	return s
}

func GetPreferenceByUserIdAndId(userId uint, id int) (Standard, bool) {
	var s Standard
	db.DB.Where("user_id = ? and id = ?", userId, id).First(&s)
	if s.ID == 0 {
		return s, false
	}
	return s, true
}

func CreateWithDefault(userID uint) Standard {
	var s Standard
	s.UserID = userID
	s.ViewTime = 10
	s.SoundVolume = 100
	s.Create()
	return s
}

type StandardDonate struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	// StandardID uint   `json:"standard_id"`
	Nickname  string     `json:"nickname"`
	Message   string     `json:"message"`
	Money     int        `json:"money"`
	Viewed    bool       `json:"viewed"`
	TransID   uint       `json:"trans_id"`
	FromOwner bool       `json:"from_owner"`
	UpdatedAt *time.Time `json:"-"`
	CreatedAt *time.Time `json:"-"`
}

func (sd *StandardDonate) Create() {
	db.DB.Create(sd)
	if sd.Viewed == false && sd.FromOwner {
		b, _ := json.Marshal(sd)
		sendAddStandardDonateWS(sd.UserID, b)
	}
}

func (sd *StandardDonate) Save() {
	db.DB.Save(sd)
	if (sd.TransID != 0 && !sd.Viewed) || sd.FromOwner {
		b, _ := json.Marshal(sd)
		sendSaveStandardDonateWS(sd.UserID, b)
	}
}

func GetDonatesByUser(userID uint) []StandardDonate {
	var donates []StandardDonate
	db.DB.Where("user_id = ? and viewed = ?", userID, false).Order("created_at").Limit(1000).Find(&donates)
	return donates
}

func GetDonateByUserAndId(userID uint, id uint) (StandardDonate, bool) {
	var d StandardDonate
	db.DB.Where("user_id = ? and id = ?", userID, id).First(&d)
	if d.ID == 0 {
		return d, false
	}
	return d, true
}

func CreateStandardDonate(nickname string, userID uint, message string, money int, fromOwner bool) StandardDonate {
	sd := StandardDonate{Nickname: nickname, UserID: userID, Message: message, FromOwner: fromOwner, Money: money}
	sd.Create()
	return sd
}

func CreateFromDonate(id uint, userID uint, money int, transID uint) error {
	// var sd StandardDonate
	// json.Unmarshal([]byte(data), &sd)
	// sd.Create()
	donate, exist := GetDonateByUserAndId(userID, id)
	if !exist {
		return errors.New("No standard donate")
	}
	donate.Money = money
	donate.TransID = transID
	donate.Save()
	return nil
}
