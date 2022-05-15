package raftServer

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/letterbeezps/raftCache/internal/cache"
	"go.uber.org/zap"
)

const (
	retainSnapShotCount = 2
	raftTimeOut         = 10 * time.Second
)

type command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value []byte `json:"value,omitempty"`
}

type RaftServer struct {
	RaftDir  string
	RaftBind string
	mu       sync.Mutex
	cache    *cache.Cache

	raft *raft.Raft
}

func NewRaftServer() *RaftServer {
	return &RaftServer{
		cache: cache.NewCache(),
	}
}

func (rs *RaftServer) StartRaftServer(enableSingle bool, localID string) error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	addr, err := net.ResolveTCPAddr("tcp", rs.RaftBind)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(rs.RaftBind, addr, 3, 10*time.Second, os.Stderr)

	if err != nil {
		return err
	}

	snapsShot, err := raft.NewFileSnapshotStore(rs.RaftDir, retainSnapShotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapShot store: %v", err)
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(rs.RaftDir, "raft-log.db"))
	if err != nil {
		return fmt.Errorf("new bolt store: %v", err)
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(rs.RaftDir, "raft-stable.db"))
	if err != nil {
		return fmt.Errorf("new bolt store: %v", err)
	}

	ra, err := raft.NewRaft(config, rs, logStore, stableStore, snapsShot, transport)

	rs.raft = ra

	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}

	return nil

}

func (rs *RaftServer) Get(key string) ([]byte, error) {
	return rs.cache.Get(key)
}

func (rs *RaftServer) Set(key string, value []byte) error {
	if rs.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	cmd := &command{
		Op:    "set",
		Key:   key,
		Value: value,
	}

	byteCmd, err := json.Marshal(cmd)
	if err != nil {
		return nil
	}

	f := rs.raft.Apply(byteCmd, raftTimeOut)
	return f.Error()
}

func (rs *RaftServer) Delete(key string) error {
	if rs.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	cmd := &command{
		Op:  "delete",
		Key: key,
	}

	byteCmd, err := json.Marshal(cmd)
	if err != nil {
		return nil
	}

	f := rs.raft.Apply(byteCmd, raftTimeOut)
	return f.Error()
}

func (rs *RaftServer) Join(nodeID, addr string) error {
	zap.S().Infof("received join request for remoted node %s at %s", nodeID, addr)

	configFuture := rs.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		zap.S().Infof("failed to get rafte configuration: %v", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			if srv.ID == raft.ServerID(nodeID) && srv.Address == raft.ServerAddress(addr) {
				zap.S().Infof("node %s at %s already exits.", nodeID, addr)
				return nil
			}

			future := rs.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}

	}
	f := rs.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	zap.S().Infof("node %s at %s joined successfully", nodeID, addr)
	return nil
}
