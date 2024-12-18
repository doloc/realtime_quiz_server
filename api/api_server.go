package api

import (
	"realtime_quiz_server/api/middleware"
	"realtime_quiz_server/api/router"
	"realtime_quiz_server/configuration"
	"realtime_quiz_server/internal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	store  *gorm.DB
	router *gin.Engine
	hub    *internal.Hub
}

func NewServer(store *gorm.DB, cf *configuration.Config) *Server {
	server := &Server{
		store: store,
		hub:   internal.NewHub(),
	}
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	const BasePath = ""
	apiGroup := r.Group(BasePath)
	{
		router.AuthRouters(apiGroup, store, cf)
		router.RoomRouters(apiGroup, store, cf)
		router.QuizRouters(apiGroup, store, cf)
		router.WebSocketRouter(apiGroup, server.hub)
	}

	r.Static("/static", "/resources")

	server.router = r

	go server.hub.Run()

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
