package task

import (
	"math/rand"
	"github.com/hashicorp/raft"
)


type Task struct {
	raftNode			*raft.Raft
	ID       			int	
	Title    			string
	FunctionToExecute 	func() error
}


func NewTask(function func() error) *Task {
	return &Task{
		ID: rand.Int(),
		Title: "Task",
		FunctionToExecute: function,
	}
}
