package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/normegil/zookeeper-rest/log"
)

const LOG_PATH string = "/tmp/"
const PORT int = 8080

var LOG *logrus.Entry = log.New(LOG_PATH, "zookeeper-rest")

func main() {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Test")
	})

	LOG.WithField("port", PORT).Info("Launch server")
	http.ListenAndServe(":"+strconv.Itoa(PORT), router)
}
