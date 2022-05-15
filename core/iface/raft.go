package iface

type RaftServer interface {
	Cache

	Join(nodeId, addr string) error
}
