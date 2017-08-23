package zookeeper

import "github.com/samuel/go-zookeeper/zk"

type ZookeeperClient interface {
	Get(string) ([]byte, Stat, error)
	Children(path string) ([]string, Stat, error)
	Delete(path string, version int32) error
	Exists(path string) (bool, Stat, error)
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	Set(path string, data []byte, version int32) (Stat, error)
	Close()
}

type Stat interface {
	Version() int32
}
