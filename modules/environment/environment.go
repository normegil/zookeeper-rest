package environment

import (
	"github.com/Sirupsen/logrus"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
)

type Env struct {
	Logger *logrus.Entry
	Zk     zookeeper.Zookeeper
}

func (e Env) Log() *logrus.Entry {
	return e.Logger
}

func (e Env) Zookeeper() zookeeper.Zookeeper {
	return e.Zk
}
