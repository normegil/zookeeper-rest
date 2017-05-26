package environment

import (
	"github.com/Sirupsen/logrus"
	"github.com/normegil/zookeeper-rest/modules/database/mongo"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
)

type Env struct {
	Logger *logrus.Entry
	Zk     zookeeper.Zookeeper
	Mongo  *mongo.Mongo
}

func (e Env) Log() *logrus.Entry {
	return e.Logger
}

func (e Env) Zookeeper() zookeeper.Zookeeper {
	return e.Zk
}

func (e *Env) Session() *mongo.Mongo {
	return e.Mongo
}
