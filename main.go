package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/normegil/zookeeper-rest/log"
	"github.com/normegil/zookeeper-rest/model"
	"github.com/normegil/zookeeper-rest/rest"
	"github.com/normegil/zookeeper-rest/router"
)

const PORT int = 8080
const LOG_PATH string = "/tmp/"

var LOG *logrus.Entry = log.New(LOG_PATH, "zookeeper-rest")

func main() {
	LOG.Logger.Level = logrus.DebugLevel
	env := model.Env{LOG}
	rt := router.New(env)
	if err := rt.Register(rest.Controller{env}.Routes()); nil != err {
		panic(err)
	}
	if err := rt.Listen(PORT); nil != err {
		panic(err)
	}
}
