package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/normegil/zookeeper-rest/api/node"
	"github.com/normegil/zookeeper-rest/modules/environment"
	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/normegil/zookeeper-rest/router"
)

const PORT int = 8080
const LOG_PATH string = "/tmp/"

var LOG *logrus.Entry = log.New(LOG_PATH, "zookeeper-rest")

func main() {
	env := environment.Env{LOG}
	rt := router.New(env)
	if err := rt.Register(node.Controller{env}.Routes()); nil != err {
		panic(err)
	}
	if err := rt.Listen(PORT); nil != err {
		panic(err)
	}
}
