package zookeeper

import (
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Zookeeper struct {
	Address            net.TCPAddr
	Logger             logrus.FieldLogger
	SubstitutionClient ZookeeperClient
}

func (z Zookeeper) Log() logrus.FieldLogger {
	return z.Logger
}

func (z Zookeeper) client() (ZookeeperClient, error) {
	if nil != z.SubstitutionClient {
		return z.SubstitutionClient, nil
	}
	if "" == z.Address.String() {
		return nil, errors.New("Empty Address/No Fake Client")
	}
	return newClient([]string{z.Address.String()}, nil, time.Second)
}
