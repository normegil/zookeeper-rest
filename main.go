package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/normegil/zookeeper-rest/log"
	"github.com/normegil/zookeeper-rest/router"
)

const PORT int = 8080
const LOG_PATH string = "/tmp/"

var LOG *logrus.Entry = log.New(LOG_PATH, "zookeeper-rest")

func main() {
	rt := &router.Router{LOG}
	rt.Serve(PORT)
}
