package db

import (
	"github.com/tarantool/go-tarantool"
	"time"
	"log"
	"os"
)

var TDB *tarantool.Connection

func tarantoolConnect() {
	host := os.Getenv("TARANTOOL_SERVICE")
	port := os.Getenv("TARANTOOL_PORT")
	if host == "" || port == "" {
		panic("please, write port ant host of tarantool")
	}
	username := os.Getenv("TARANTOOL_USER_NAME")
	password := os.Getenv("TARANTOOL_USER_PASSWORD")
	server := host + ":" + port
	opts := tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          username,
		Pass:          password,
	}
	client, err := tarantool.Connect(server, opts)
	if err != nil {
		log.Fatal(err)
	}
	TDB = client
}
