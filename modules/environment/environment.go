package environment

import (
	"github.com/Sirupsen/logrus"
)

type Env struct {
	Logger *logrus.Entry
}

func (e Env) Log() *logrus.Entry {
	return e.Logger
}
