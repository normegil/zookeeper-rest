package zookeeper

import (
	"time"

	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/samuel/go-zookeeper/zk"
)

func (z Zookeeper) Create(path string, content []byte, acls []zk.ACL) error {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return err
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	var noFlag int32 = 0
	_, err = connection.Create(path, content, noFlag, acls)
	if nil != err {
		return err
	}
	return nil
}
