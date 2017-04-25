package zookeeper

import (
	"time"

	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

func (z Zookeeper) Delete(path string) error {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return errors.Wrap(err, "Connecting to Zookeeper")
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	_, stat, err := connection.Get(path)
	if nil != err {
		return errors.Wrapf(err, "Loading statistics of %s", path)
	}
	return errors.Wrapf(connection.Delete(path, stat.Version), "Deleting %s", path)
}
