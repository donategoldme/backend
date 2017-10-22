package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB = postgresConnect()

func postgresConnect() *gorm.DB {
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
	return db
}
