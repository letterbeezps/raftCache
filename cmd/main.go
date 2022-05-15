package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/letterbeezps/raftCache/core/util"
	"github.com/letterbeezps/raftCache/global"
	"github.com/letterbeezps/raftCache/initial"
	"go.uber.org/zap"
)

const (
	Host            = "127.0.0.1"
	DefaultHTTPAddr = ":11000"
	DefaultRaftAddr = ":12000"
	RaftDirPrefix   = "./test_db"
)

var (
	httpAddr = flag.String("haddr", DefaultHTTPAddr, "set the deafult http bind addr")
	raftAddr = flag.String("raddr", DefaultRaftAddr, "set the raft bind addr")
	joinAddr = flag.String("join", "", "set the join addr")
	nodeID   = flag.String("id", "", "Node Id")
)

func initOption() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No raft storage dir specified\n")
		os.Exit(1)
	}
}

func main() {
	initOption()

	initial.InitLogger()

	initial.InitRaftServer()

	raftDir := RaftDirPrefix + flag.Arg(0)
	if raftDir == RaftDirPrefix {
		fmt.Fprintf(os.Stderr, "No raft storage dir specified\n")
		os.Exit(1)
	}

	global.RaftServer.RaftBind = Host + *raftAddr
	global.RaftServer.RaftDir = raftDir

	if err := global.RaftServer.StartRaftServer(*joinAddr == "", *nodeID); err != nil {
		zap.S().Fatalf("failed to start Raft Server: %s", err.Error())
		panic(fmt.Sprintf("failed to start Raft Server: %s", err.Error()))
	}

	if *joinAddr != "" {
		if err := util.JoinRaft(Host+*joinAddr, Host+*raftAddr, *nodeID); err != nil {
			zap.S().Fatalf("failed to Join Node %s at %s", Host+*joinAddr, err.Error())
			panic(fmt.Sprintf("failed to Join Node %s at %s", Host+*joinAddr, err.Error()))
		}
	}

	Router := initial.Routers()

	zap.S().Info("start server at: %s", fmt.Sprintf("%s", *httpAddr))

	if err := Router.Run(fmt.Sprintf("%s", *httpAddr)); err != nil {
		zap.S().Panic("start server failed", err.Error())
	}
}
