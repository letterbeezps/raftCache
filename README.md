# raftCache

一个用来学习raft协议的demo，使用raft协议来构建一个分布式缓存，项目参考了[hraftd](https://github.com/otoolep/hraftd)。

## 调试代码

**leader** 启动初始化节点node0

```shell
cd /cmd
go run main.go -id node0 node0
```

**follower** 加入多个follower node1 & node2

```shell
go run main.go -id node1 -haddr :11001 -raddr :12001 -join :11000 node1

go run main.go -id node2 -haddr :11002 -raddr :12002 -join :11000 node2
```

得到三个副本

```shell
├── test_dbnode0
│   ├── raft-log.db
│   ├── raft-stable.db
│   └── snapshots
├── test_dbnode1
│   ├── raft-log.db
│   ├── raft-stable.db
│   └── snapshots
└── test_dbnode2
    ├── raft-log.db
    ├── raft-stable.db
    └── snapshots
```

## 调试API

只能对leader节点发起**写**请求

> POST
> 127.0.0.1:11000/v1/cache/key

```json
{
    "value": "testNode0Value1"
}
```

> DELETE
> 127.0.0.1:11000/v1/cache/key

可以最任意一个节点发起**读**请求

> GET
> 127.0.0.1:11000/v1/cache/key
> 127.0.0.1:11001/v1/cache/key
> 127.0.0.1:11002/v1/cache/key

```json
{
    "key": "key",
    "value": {
        "value": "testNode0Value1"
    }
}
```
