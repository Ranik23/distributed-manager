package main

import (
	"distributed-manager/server/internal/api/handlers"
	"distributed-manager/server/internal/db"
	"distributed-manager/server/internal/entities/task"
	"distributed-manager/server/internal/server"
	"distributed-manager/server/internal/server/config"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



func main() {

	Logger := slog.Default()

	Config := config.NewConfig("localhost", "8080")

	dsn := "host=localhost user=postgres password=postgres dbname=test_new port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}


	database := storage.NewDatabase(db, Logger)

	database.AutoMigrate(&task.Task{})


	r := gin.Default()

	server := server.NewServer(database, Config, r)

	r.POST("/tasks", handlers.TaskCreate(server))
	r.GET("/tasks/:id", handlers.GetTask(*server))

	if err := server.Run(); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }

}