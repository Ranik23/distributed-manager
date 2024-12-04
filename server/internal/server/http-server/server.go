package server

import (
	storage "distributed-manager/server/internal/db/postgres"
	"log/slog"

	"distributed-manager/server/internal/server/grpc"

	"distributed-manager/server/internal/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	DB           storage.Storage
	Config *config.Config
	Router       *gin.Engine
	GRPCServer   *grpc.RaftGRPCServiceServer
	Logger       *slog.Logger
}

func NewServer(db storage.Storage,
	Config *config.Config,
	router *gin.Engine,
	grpcServer *grpc.RaftGRPCServiceServer,
	logger *slog.Logger) *Server {

	return &Server{
		DB:           db,
		Config: 	  Config,
		Router:       router,
		GRPCServer:   grpcServer,
		Logger:       logger,
	}
}

func (s *Server) Run() error {
	return s.Router.Run(s.Config.Server.HTTP.Host+ ":" + s.Config.Server.HTTP.Port)
}
