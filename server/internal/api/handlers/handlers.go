package handlers

import (
	"distributed-manager/server/internal/server"
	"github.com/gin-gonic/gin"
)



func TaskCreate(srv *server.Server) func (*gin.Context) {
	// empty
	return func(g *gin.Context) {
		g.JSON(201, gin.H{"message": "Task Created"})
	}
}

func GetTask(srv server.Server) func (*gin.Context) {
	//empty
	return func(g *gin.Context) {
		taskID := g.Param("id")
		g.JSON(200, gin.H{"taskID": taskID, "status": "completed"})
	}
}