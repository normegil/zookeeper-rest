package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/normegil/zookeeper-rest/api/node"
	"github.com/normegil/zookeeper-rest/modules/environment"
	"github.com/normegil/zookeeper-rest/modules/log"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
	"github.com/normegil/zookeeper-rest/router"
	"github.com/pkg/errors"
)

const PORT int = 8080
const LOG_PATH string = "/tmp/"
const ZK_ADDRESS string = "127.0.0.1"
const VERBOSE bool = true

var LOG *logrus.Entry = log.New(LOG_PATH, "zookeeper-rest", VERBOSE)

func main() {
	env := environment.Env{LOG, zookeeper.Zookeeper{ZK_ADDRESS, LOG}}
	rt := router.New(env)
	if err := rt.Register(node.Controller{env}.Routes()); nil != err {
		panic(errors.Wrap(err, "Could not register Node controllers: "))
	}
	if err := rt.Listen(PORT); nil != err {
		panic(errors.Wrapf(err, "Fatal error while server listening (port:%d)", PORT))
	}
}
