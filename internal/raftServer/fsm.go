package raftServer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/raft"
	"github.com/letterbeezps/raftCache/internal/cache"
	"go.uber.org/zap"
)

// type fsm RaftServer

func (fsm *RaftServer) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		zap.S().Infof("failed to unmarshal command: %s", err.Error())
		// panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case "set":
		return fsm.applySet(c.Key, c.Value)
	case "delete":
		return fsm.applyDelete(c.Key)
	default:
		zap.S().Infof("unrecognized ccommand op: %s", c.Op)
		panic(fmt.Sprintf("unrecognized ccommand op: %s", c.Op))
	}
}

func (fsm *RaftServer) applySet(key string, value []byte) interface{} {
	return fsm.cache.Set(key, value)
}

func (fsm *RaftServer) applyDelete(key string) interface{} {
	return fsm.cache.Delete(key)
}

func (fsm *RaftServer) Snapshot() (raft.FSMSnapshot, error) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	byteCache, err := fsm.cache.Dump()
	if err != nil {
		panic(fmt.Sprintf("cant dump cache %s", err.Error()))
	}

	return &fsmSnapShot{data: byteCache}, nil
}

func (fsm *RaftServer) Restore(rc io.ReadCloser) error {
	db := make(map[string][]byte)
	if err := json.NewDecoder(rc).Decode(&db); err != nil {
		return err
	}

	fsm.cache = cache.NewCacheByDb(db)
	return nil
}
