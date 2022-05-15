package initial

import (
	"github.com/letterbeezps/raftCache/global"
	"github.com/letterbeezps/raftCache/internal/raftServer"
)

func InitRaftServer() {
	global.RaftServer = raftServer.NewRaftServer()
}
