package raft

import (
	"encoding/json"
	"github.com/hashicorp/raft"
)

// RaftStoreSnapshot структура для снимка состояния FSM
type RaftStoreSnapshot struct {
	tasks []Task
}

// Persist сохраняет снимок
func (s *RaftStoreSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		enc := json.NewEncoder(sink)
		if err := enc.Encode(s.tasks); err != nil {
			return err
		}
		return nil
	}()

	if err != nil {
		sink.Cancel()
		return err
	}

	return sink.Close()
}

// Release освобождает ресурсы снимка
func (s *RaftStoreSnapshot) Release() {}