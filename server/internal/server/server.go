package server

import (
	storage "distributed-manager/server/internal/db/postgres"
	"distributed-manager/server/internal/entities/manager"
	//"distributed-manager/server/internal/models"
	"distributed-manager/server/internal/server/config"

	"github.com/gin-gonic/gin"
)




type Server struct {
	DB storage.Storage
	ServerConfig *config.Config
	Router *gin.Engine
	Manager *manager.Manager
}


func NewServer(db storage.Storage, serverConfig *config.Config, router *gin.Engine) *Server {
	return &Server{
		DB: db,
		ServerConfig: serverConfig, 
		Router: router,
	}
}

func (s *Server) Run() error {
	return s.Router.Run(s.ServerConfig.Host + ":" + s.ServerConfig.Port)
}