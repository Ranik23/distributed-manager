package handlers

import (
	"context"
	"distributed-manager/server/internal/server/http-server"
	raftservice "distributed-manager/server/pkg"
	"github.com/gin-gonic/gin"
)



func TaskCreateHandler(srv *server.Server) func (*gin.Context) {
	return func(g *gin.Context) {
		request := raftservice.TaskRequest{
			Task: "Hello",
		}
		response, err := srv.GRPCServer.SubmitTask(context.Background(), &request)
		if err != nil {
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
		g.JSON(201, gin.H{"message": "Task Created", "success": response.Success})
	}
}

func TaskGetHandler(srv *server.Server) func (*gin.Context) {
	return func(g *gin.Context) {
		taskID := g.Param("id")

		request := raftservice.GetTaskRequest{
			Id: taskID,
		}

		response, err := srv.GRPCServer.GetTask(context.Background(), &request)
		if err != nil {
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}

		g.JSON(200, gin.H{"taskID": taskID, "status": "completed", "content": response.GetData()})
	}
}
