package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

const PORT int = 8080

var log *logrus.Entry = logrus.NewEntry(logrus.New())

func init() {
	log = log.WithField("Test", "123")
	log.Logger.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.InfoLevel:  "info.log",
		logrus.ErrorLevel: "error.log",
	}))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test")
	})

	log.WithField("port", PORT).Info("Launch server")
	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
}
