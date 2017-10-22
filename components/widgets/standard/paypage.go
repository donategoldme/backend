package standard

import (
	"encoding/json"
	"fmt"
	"time"

	"donategold.me/db"
)

func init() {
	db.DB.AutoMigrate(&Paypage{})
}

type Paypage struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`

	HeaderTitle string `json:"header_title"`
	HeaderDesc  string `json:"header_desc"`
	//здесь хранится array уссылок на каналы
	HeaderChans     string `json:"header_chans"`
	HeaderTextColor string `json:"header_text_color"`
	HeaderImg       string `json:"header_img"`
	HeaderColor     string `json:"header_color"`
	HeaderFont      string `json:"header_font"`

	BackgroundColor string `json:"background_color"`
	BackgroundImg   string `json:"background_img"`

	MainTitle      string `json:"main_title"`
	MainTitleColor string `json:"main_title_color"`
	MainDesc       string `json:"main_desc"`
	MainDescColor  string `json:"main_desc_color"`
	MainFont       string `json:"main_font"`
	ButtonColor    string `json:"button_color"`
	ButtonText     string `json:"button_text"`

	DeletedAt *time.Time `json:"-"`
}

func (p *Paypage) Create() {
	db.DB.Create(p)
	data, err := json.Marshal(p)
	if err == nil {
		sendCreatePaypageWS(p.UserID, data)
	}
}

func (p *Paypage) Save() {
	if p.Active == true {
		tx := db.DB.Begin()
		tx.Model(&Paypage{}).Where("user_id = ?", p.UserID).Update("active", false)
		tx.Save(p)
		tx.Commit()
	} else {
		db.DB.Save(p)
	}
	data, err := json.Marshal(p)
	if err == nil {
		sendSavePaypageWS(p.UserID, data)
	}
}

func GetByUserID(userID uint) []Paypage {
	var p []Paypage
	db.DB.Where("user_id = ?", userID).Order("id").Find(&p)
	return p
}

func GetByID(id uint) (p Paypage, exist bool) {
	db.DB.Where("id = ?", id).First(&p)
	if p.ID != 0 {
		exist = true
	}
	return
}

func GetByUserIDAndID(userID uint, id uint) (p Paypage, exist bool) {
	db.DB.Where("id = ? and user_id = ?", id, userID).First(&p)
	if p.ID != 0 {
		exist = true
	}
	return
}

func DeleteByUserIDAndID(userID uint, id uint) {
	db.DB.Where("user_id = ? and id = ?", userID, id).Delete(&Paypage{})
	sendDeletePaypageWS(userID, []byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func GetActivePaypageByUserID(userID uint) (p Paypage, exist bool) {
	db.DB.Where("user_id = ? and active = ?", userID, true).First(&p)
	if p.ID != 0 {
		exist = true
	}
	return
}
