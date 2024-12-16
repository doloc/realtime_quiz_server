package main

import (
	"log"
	"net/http"
	"realtime_quiz_server/api"
	"realtime_quiz_server/cache"
	"realtime_quiz_server/configuration"
	"realtime_quiz_server/database"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	config, err := configuration.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration: ", err)
	}

	db, err := database.OpenConnectionToDatabase(config.DSN)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	// Khởi tạo kết nối Redis
	cache.InitializeRedis(config.RedisAddress, config.RedisPassword, config.RedisDB)

	server := api.NewServer(db, &config)

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
