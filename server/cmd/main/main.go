package main

import (
	"distributed-manager/server/internal/api/handlers"
	boltDb "distributed-manager/server/internal/db/bolt"
	"os"
	"path/filepath"

	"distributed-manager/server/internal/raft"
	"distributed-manager/server/internal/server/grpc"
	"distributed-manager/server/internal/server/http-server"
	"distributed-manager/server/internal/config"
	"log"
	"log/slog"
	"github.com/gin-gonic/gin"
)



func main() {

	Logger := slog.Default()

	Config, err := config.LoadConfig(filepath.Join(os.Getenv("HOME"), "disributed-manager/server/internal/config/config.yaml"))
	if err != nil {
		log.Fatalf("error")
	}

	raftLogPath := "disributed-manager/trash/raft-log.db"

	bolt_db, err := boltDb.NewBoltDatabase(filepath.Join(os.Getenv("HOME"), raftLogPath))//filepath.Join(os.Getenv("HOME"), raftLogPath))
	if err != nil {
		log.Fatalf("Failed to create bolr db: %v", err)
	}


	raft, err := raft.NewRaftStore(filepath.Join(os.Getenv("HOME"), raftLogPath), "localhost:9091")
	if err != nil {
		log.Fatal(err)
	}


	raftGRPCServiceServer := grpc.StartGRPCRaftServer(raft)


	r := gin.Default()
	server := server.NewServer(bolt_db, Config, r, raftGRPCServiceServer, Logger)
	r.POST("/tasks", handlers.TaskCreateHandler(server))
	r.GET("/tasks/:id", handlers.TaskGetHandler(server))

	if err := server.Run(); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }

}