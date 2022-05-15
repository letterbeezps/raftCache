package raftServer

import (
	"github.com/hashicorp/raft"
)

// 定义一个快照结构
type fsmSnapShot struct {
	data []byte
}

func (f *fsmSnapShot) Persist(sink raft.SnapshotSink) error {
	err := func() error {

		if _, err := sink.Write(f.data); err != nil {
			return err
		}

		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapShot) Release() {}
