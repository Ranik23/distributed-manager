package raft

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

// Task структура для хранения задачи
type Task struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// RaftStore структура для хранения данных в Raft
type RaftStore struct {
	raftDir  string
	raftBind string
	raft     *raft.Raft
	tasks    []Task
}

// Apply применяет команду Raft к FSM
func (s *RaftStore) Apply(l *raft.Log) interface{} {
	var task Task
	if err := json.Unmarshal(l.Data, &task); err != nil {
		log.Printf("failed to unmarshal task: %s", err.Error())
		return err
	}

	s.tasks = append(s.tasks, task)
	return nil
}

// Snapshot возвращает текущее состояние FSM
func (s *RaftStore) Snapshot() (raft.FSMSnapshot, error) {
	return &RaftStoreSnapshot{tasks: s.tasks}, nil
}

// Restore восстанавливает состояние FSM из снимка
func (s *RaftStore) Restore(rc io.ReadCloser) error {
	var tasks []Task
	if err := json.NewDecoder(rc).Decode(&tasks); err != nil {
		return err
	}

	s.tasks = tasks
	return nil
}


// NewRaftStore создаёт новый экземпляр RaftStore
func NewRaftStore(raftDir, raftBind string) (*RaftStore, error) {
	s := &RaftStore{
		raftDir:  raftDir,
		raftBind: raftBind,
		tasks:    make([]Task, 0),
	}

	// Создание директории для Raft
	if err := os.MkdirAll(raftDir, 0755); err != nil {
		return nil, err
	}

	// Инициализация хранилища логов
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-log.db"))
	if err != nil {
		return nil, err
	}

	// Инициализация хранилища снимков
	snapshotStore, err := raft.NewFileSnapshotStore(raftDir, 1, os.Stderr)
	if err != nil {
		return nil, err
	}

	// Инициализация транспортного слоя
	addr, err := net.ResolveTCPAddr("tcp", raftBind)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(raftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	// Инициализация конфигурации Raft
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftBind)

	// Инициализация Raft
	ra, err := raft.NewRaft(config, s, logStore, logStore, snapshotStore, transport)
	if err != nil {
		return nil, err
	}

	s.raft = ra

	// Если узел первый, инициализировать кластер
	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(raftBind),
				Address: raft.ServerAddress(raftBind),
			},
		},
	}

	if err := s.raft.BootstrapCluster(configuration).Error(); err != nil && err != raft.ErrCantBootstrap {
		return nil, err
	}

	// Проверка роли узла
	if s.raft.State() != raft.Leader {
		log.Println("This node is not the leader")
	} else {
		log.Println("This node is the leader")
	}

	return s, nil
}

// AddNode добавляет новый узел в кластер
func (s *RaftStore) AddNode(newNodeID, newNodeAddress string) error {
	if s.raft.State() != raft.Leader {
		return raft.ErrNotLeader
	}

	if err := s.raft.AddVoter(raft.ServerID(newNodeID), raft.ServerAddress(newNodeAddress), 0, 0).Error(); err != nil {
		log.Printf("Failed to add node: %s", err.Error())
		return err
	}

	log.Printf("Node %s added successfully", newNodeID)
	return nil
}
