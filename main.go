package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	logrotation "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
)

const LOG_PATH string = "/tmp/"
const PORT int = 8080

var log *logrus.Entry = logrus.NewEntry(logrus.New())

func init() {
	extention := "log"
	hook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  NewLogRotation(LOG_PATH, "zookeeper-rest.info", extention),
		logrus.ErrorLevel: NewLogRotation(LOG_PATH, "zookeeper-rest.error", extention),
	})
	hook.SetFormatter(&logrus.JSONFormatter{})
	log.Logger.Hooks.Add(hook)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test")
	})

	log.WithField("port", PORT).Info("Launch server")
	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
}

func NewLogRotation(path string, name string, extention string) *logrotation.RotateLogs {
	pattern := "%Y-%m-%d"
	separator := "."

	day := time.Duration(24) * time.Hour
	return logrotation.New(
		path+name+separator+pattern+separator+extention,
		logrotation.WithLinkName(path+name+separator+extention),
		logrotation.WithMaxAge(time.Duration(7)*day),
	)
}
