package zookeeper

import (
	"github.com/Sirupsen/logrus"
)

type Zookeeper struct {
	Address string
	Logger  *logrus.Entry
}

func (z Zookeeper) Log() *logrus.Entry {
	return z.Logger
}

type Node struct {
	Path    string
	Content string
	Childs  []string
}
