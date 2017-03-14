package zookeeper

import (
	"time"

	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/samuel/go-zookeeper/zk"
)

func (z Zookeeper) Set(path string, content []byte, version int) error {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return err
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	_, err = connection.Set(path, content, int32(version))
	if nil != err {
		return  err
	}
	return nil
}
