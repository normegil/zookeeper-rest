package zookeeper

import (
	"time"

	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

func (z Zookeeper) Children(path string) ([]string, error) {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return nil, errors.Wrap(err, "Connecting to Zookeeper")
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	z.Log().WithField("parentPath", path).Debug("Load childs")
	children, _, err := connection.Children(path)
	if nil != err {
		return nil, errors.Wrapf(err, "Loading childs for %s", path)
	}
	return children, nil
}
