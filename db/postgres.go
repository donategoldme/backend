package db

import (
	_ "github.com/lib/pq"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
)

var DB *gorm.DB

func postgresConnect() {
	service := os.Getenv("DB_SERVICE")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	port := os.Getenv("DB_PORT")
	params := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", user, name, pass, service,
	port)
	db, err := gorm.Open("postgres", params)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	DB = db
}