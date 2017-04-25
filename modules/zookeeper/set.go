package zookeeper

import (
	"time"

	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

func (z Zookeeper) Set(path string, content []byte) error {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return errors.Wrap(err, "Connecting to zookeeper")
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	exist, _, err := connection.Exists(path)
	if nil != err {
		return errors.Wrap(err, "Testing path for existence")
	}

	z.Log().WithField("path", path).WithField("Creation", !exist).Debug("Node Creation/Update")
	if exist {
		_, stat, err := connection.Get(path)
		if nil != err {
			return errors.Wrap(err, "Getting path version")
		}
		_, err = connection.Set(path, content, stat.Version)
	} else {
		var noFlag int32 = 0
		_, err = connection.Create(path, content, noFlag, zk.WorldACL(zk.PermAll))
		err = errors.Wrap(err, "Creating path")
	}

	return err
}
