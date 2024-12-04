package grpc

import (
	"context"
	"distributed-manager/server/internal/raft"
	raftservice "distributed-manager/server/pkg"
	"log"
	"net"

	"google.golang.org/grpc"
)

type RaftGRPCServiceServer struct {
	raftservice.UnimplementedTaskServiceServer
	store *raft.RaftStore
}

func (s *RaftGRPCServiceServer) SubmitTask(ctx context.Context, req *raftservice.TaskRequest) (*raftservice.TaskResponse, error) {
	// TODO
	return &raftservice.TaskResponse{Success: true}, nil
}

func (s *RaftGRPCServiceServer) GetTask(ctx context.Context, req *raftservice.GetTaskRequest) (*raftservice.GetTaskResponse, error) {
	// TODO
	return &raftservice.GetTaskResponse{
		Data: "Hello",
	}, nil
}


func StartGRPCRaftServer(store *raft.RaftStore) *RaftGRPCServiceServer {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10 * 1024 * 1024),  
		grpc.MaxSendMsgSize(10 * 1024 * 1024),
	)


	raftServiceServer := RaftGRPCServiceServer{store: store}
	
	raftservice.RegisterTaskServiceServer(grpcServer, &raftServiceServer)

	log.Println("Starting gRPC server at port 50051")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return &raftServiceServer
}
