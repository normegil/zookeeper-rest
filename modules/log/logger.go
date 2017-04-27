package log

import (
	"strings"
	"time"

	stackhook "github.com/Gurpartap/logrus-stack"
	"github.com/Sirupsen/logrus"
	logrotation "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
)

func New(path string, filename string, verbose bool) *logrus.Entry {
	extention := "log"
	fileHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  newLogRotation(path, filename+".info", extention),
		logrus.ErrorLevel: newLogRotation(path, filename+".error", extention),
	})
	fileHook.SetFormatter(&logrus.JSONFormatter{})

	log := logrus.NewEntry(logrus.New())

	var hook stackhook.LogrusStackHook
	var logLvl logrus.Level
	if verbose {
		logLvl = logrus.DebugLevel
		stackLvl := []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		}
		hook = stackhook.NewHook(logrus.AllLevels, stackLvl)
	} else {
		logLvl = logrus.InfoLevel
		stackLvl := []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		}
		callerLvl := []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		}
		hook = stackhook.NewHook(callerLvl, stackLvl)
	}
	log.Logger.Level = logLvl
	log.Logger.Hooks.Add(hook)
	log.Logger.Hooks.Add(fileHook)

	return log
}

func newLogRotation(path string, name string, extention string) *logrotation.RotateLogs {
	pattern := "%Y-%m-%d"
	separator := "."

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	day := time.Duration(24) * time.Hour
	return logrotation.New(
		path+name+separator+pattern+separator+extention,
		logrotation.WithLinkName(path+name+separator+extention),
		logrotation.WithMaxAge(time.Duration(7)*day),
	)
}
