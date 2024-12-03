package main

import (
	"io"
	Log "log"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

type MyFSM struct{}

func (f *MyFSM) Apply(Log *raft.Log) interface{} {
	// Реализация применения команды
	return nil
}

func (f *MyFSM) Snapshot() (raft.FSMSnapshot, error) {
	// Реализация создания снимка состояния
	return nil, nil
}

func (f *MyFSM) Restore(rc io.ReadCloser) error {
	// Реализация восстановления из снимка состояния
	return nil
}

func createTransport(address string) (*raft.NetworkTransport, error) {
	return raft.NewTCPTransport(address, nil, 2, time.Second, os.Stderr)
}

func main() {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID("node1")

	LogStore, err := raftboltdb.NewBoltStore("raft.db")
	if err != nil {
		Log.Fatalf("failed to create Log store: %v", err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore("snapshots", 3, os.Stderr)
	if err != nil {
		Log.Fatalf("failed to create snapshot store: %v", err)
	}

	transport, err := createTransport("localhost:8081")
	if err != nil {
		Log.Fatalf("failed to create transport: %v", err)
	}

	fsm := &MyFSM{}
	r, err := raft.NewRaft(config, fsm, LogStore, LogStore, snapshotStore, transport)
	if err != nil {
		Log.Fatalf("failed to create Raft: %v", err)
	}

	configs := raft.Configuration{
		Servers: []raft.Server{
			{ID: "node1", Address: raft.ServerAddress("localhost:8081")},
			{ID: "node2", Address: raft.ServerAddress("localhost:8082")},
			{ID: "node3", Address: raft.ServerAddress("localhost:8083")},
		},
	}

	r.BootstrapCluster(configs)

	Log.Println("Raft node started")
}
