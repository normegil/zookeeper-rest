package zookeeper

import(
	"time"
	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/samuel/go-zookeeper/zk"
)

func (z Zookeeper) Children(path string) ([]string, error) {
	connection, _, err := zk.Connect([]string{z.Address}, time.Second)
	if nil != err {
		return nil, err
	}
	defer connection.Close()
	connection.SetLogger(log.VoidLogger{})

	children, _, err := connection.Children(path)
	if nil != err {
		return nil, err
	}
	return children, nil
}
