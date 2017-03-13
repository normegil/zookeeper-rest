package log

import (
	"time"

	"github.com/Sirupsen/logrus"
	logrotation "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
)

func New(path string, filename string) *logrus.Entry {
	extention := "log"
	hook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  newLogRotation(path, filename + ".info", extention),
		logrus.ErrorLevel: newLogRotation(path, filename + ".error", extention),
	})
	hook.SetFormatter(&logrus.JSONFormatter{})

	log := logrus.NewEntry(logrus.New())
	log.Logger.Hooks.Add(hook)
	return log
}

func newLogRotation(path string, name string, extention string) *logrotation.RotateLogs {
	pattern := "%Y-%m-%d"
	separator := "."

	day := time.Duration(24) * time.Hour
	return logrotation.New(
		path+name+separator+pattern+separator+extention,
		logrotation.WithLinkName(path+name+separator+extention),
		logrotation.WithMaxAge(time.Duration(7)*day),
	)
}
