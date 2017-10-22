package money

import (
    "github.com/jinzhu/gorm"
    "time"
    "donategold.me/db"
)

func init() {
    db.DB.AutoMigrate(&Balance{})
}

type Balance struct {
    ID uint `json:"-"`
    UserID uint `gorm:"unique_index" json:"-"`
    Money int `json:"money"`
    UpdatedAt time.Time `json:"updated_at"`
}

func(b *Balance) Add(sum int) {
    err := db.DB.Model(b).UpdateColumn("money", gorm.Expr("money + ?", sum)).Error
    if err == nil {
        b.Money += sum
    }
}

func AddTrans(tx *gorm.DB, money int, user_id uint) {
    tx.Model(&Balance{}).Where("user_id = ?", user_id).UpdateColumn("money", gorm.Expr("money + ?", money))
}

func GetBalanceByUserId(user_id uint) Balance {
    var balance Balance
    db.DB.Where("user_id = ?", user_id).First(&balance)
    if balance.ID == 0 && user_id != 0 {
        balance.UserID = user_id
        db.DB.Create(&balance)
    }
    return balance
}


