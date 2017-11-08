package environment

import (
	"github.com/normegil/mongo"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
	"github.com/sirupsen/logrus"
)

type Env struct {
	Logger logrus.FieldLogger
	Zk     zookeeper.Zookeeper
	Mongo  mongo.Session
}

func (e Env) Log() logrus.FieldLogger {
	return e.Logger
}

func (e Env) Zookeeper() zookeeper.Zookeeper {
	return e.Zk
}

func (e *Env) Session() mongo.Session {
	return e.Mongo
}
