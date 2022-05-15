package raftServer

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_DbOPen(t *testing.T) {
	rs := NewRaftServer()
	tmpDir, _ := ioutil.TempDir("", "store_test")
	fmt.Println(tmpDir)

	defer os.RemoveAll(tmpDir)

	rs.RaftBind = "127.0.0.1:0"
	rs.RaftDir = tmpDir

	if rs == nil {
		t.Fatalf("failed to create rs")
	}

	if err := rs.StartRaftServer(false, "node0"); err != nil {
		t.Fatalf("failed to open db: %s", err.Error())
	}
}

func Test_DbOPenSingleNode(t *testing.T) {
	rs := NewRaftServer()
	tmpDir, _ := ioutil.TempDir("", "store_test")
	fmt.Println(tmpDir)

	defer os.RemoveAll(tmpDir)

	rs.RaftBind = "127.0.0.1:11000"
	rs.RaftDir = tmpDir

	if rs == nil {
		t.Fatalf("failed to create rs")
	}

	if err := rs.StartRaftServer(true, "node0"); err != nil {
		t.Fatalf("failed to open db: %s", err.Error())
	}

	time.Sleep(3 * time.Second)

	if err := rs.Set("key1", []byte("value1")); err != nil {
		t.Fatalf("failed to set key: %s", err.Error())
	}

	time.Sleep(500 * time.Millisecond)
	value, err := rs.Get("key1")
	if err != nil {
		t.Fatalf("failed to get key: %s", err.Error())
	}
	if string(value) != "value1" {
		t.Fatalf("key has wrong value: %s", value)
	}

	if err := rs.Delete("key1"); err != nil {
		t.Fatalf("failed to delete key: %s", err.Error())
	}

	time.Sleep(500 * time.Millisecond)

	value, err = rs.Get("key1")
	if err != nil {
		t.Fatalf("failed to get key: %s", err.Error())
	}

	if string(value) != "" {
		t.Fatalf("key has wrong value: %s", value)
	}
}
