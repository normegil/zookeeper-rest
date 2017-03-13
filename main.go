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
	rt := router.New(LOG)
	env := model.Env{LOG}
	rt.Register(rest.Controller{env}.Routes())
	rt.Listen(PORT)
}
