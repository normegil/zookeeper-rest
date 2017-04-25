package zookeeper

import (
	"time"

	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

func (z Zookeeper) Content(path string) (string, error) {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return "", errors.Wrap(err, "Connecting to Zookeeper")
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	content, _, err := connection.Get(path)
	if nil != err {
		return "", errors.Wrapf(err, "Loading node %s", path)
	}
	return string(content), nil
}

func (z Zookeeper) Exist(path string) (bool, error) {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return false, errors.Wrap(err, "Connecting to Zookeeper")
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	exist, _, err := connection.Exists(path)
	if nil != err {
		return false, errors.Wrapf(err, "Checking existence of %s", path)
	}
	return exist, nil
}

func (z Zookeeper) Load(path string) (Node, error) {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return Node{}, errors.Wrap(err, "Connecting to Zookeeper")
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	content, _, err := connection.Get(path)
	if nil != err {
		return Node{}, errors.Wrapf(err, "Loading node %s", path)
	}

	childs, _, err := connection.Children(path)
	if nil != err {
		return Node{}, errors.Wrapf(err, "Loading childs of %s", path)
	}

	basePath := path + "/"
	if "/" == path {
		basePath = path
	}
	var childPaths []string
	for _, cPath := range childs {
		childPaths = append(childPaths, basePath+cPath)
	}

	return Node{
		Path:    path,
		Content: string(content),
		Childs:  childPaths,
	}, nil
}
