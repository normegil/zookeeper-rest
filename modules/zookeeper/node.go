package zookeeper

import (
	"time"

	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/samuel/go-zookeeper/zk"
)

func (z Zookeeper) Content(path string) (string, error) {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return "", err
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	content, _, err := connection.Get(path)
	if nil != err {
		return "", err
	}
	return string(content), nil
}

func (z Zookeeper) Exist(path string) (bool, error) {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return false, err
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	exist, _, err := connection.Exists(path)
	if nil != err {
		return false, err
	}
	return exist, nil
}
