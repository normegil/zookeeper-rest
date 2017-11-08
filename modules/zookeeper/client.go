package zookeeper

import (
	"time"

	log "github.com/normegil/golog"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

type clientImpl struct {
	*zk.Conn
}

func newClient(addresses []string, l zk.Logger, timeout time.Duration) (*clientImpl, error) {
	conn, ch, err := zk.Connect(addresses, timeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Connecting to %s", addresses)
	}
	if nil != l {
		conn.SetLogger(l)
	} else {
		conn.SetLogger(log.VoidLogger{})
	}
	for _ = range ch {
		if zk.StateHasSession == conn.State() {
			break
		}
	}
	return &clientImpl{conn}, nil
}

func (c *clientImpl) Close() {
	c.Conn.Close()
}

func (c *clientImpl) Get(path string) ([]byte, Stat, error) {
	bytes, stat, err := c.Conn.Get(path)
	return bytes, &statImpl{stat}, err
}

func (c *clientImpl) Children(path string) ([]string, Stat, error) {
	children, zkStat, err := c.Conn.Children(path)
	return children, &statImpl{zkStat}, err
}

func (c *clientImpl) Exists(path string) (bool, Stat, error) {
	exist, stat, err := c.Conn.Exists(path)
	return exist, &statImpl{stat}, err
}

func (c *clientImpl) Set(path string, data []byte, version int32) (Stat, error) {
	stat, err := c.Conn.Set(path, data, version)
	return &statImpl{stat}, err
}

type statImpl struct {
	zkStat *zk.Stat
}

func (s statImpl) Version() int32 {
	return s.zkStat.Version
}
