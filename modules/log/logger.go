package log

import (
	"strings"
	"time"

	stackhook "github.com/Gurpartap/logrus-stack"
	"github.com/Sirupsen/logrus"
	logrotation "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	uuid "github.com/satori/go.uuid"
	"github.com/weekface/mgorus"
)

type Options struct {
	Verbose bool
	File    FileOptions
	DB      MongoOptions
}

type FileOptions struct {
	FolderPath string
	FileName   string
	MaxAge     time.Duration
}

type MongoOptions struct {
	Address    string
	Port       string
	DB         string
	Collection string
	User       string
	Password   string
}

func New(opts Options) *logrus.Entry {
	log := logrus.NewEntry(logrus.New())
	if opts.Verbose {
		log.Logger.Level = logrus.DebugLevel
	}

	log.Logger.Hooks.Add(stackHK())

	if "" != opts.File.FileName {
		log.Logger.Hooks.Add(fileHK(opts.File))
	}

	if "" != opts.DB.Address {
		log.Logger.Hooks.Add(mongoHK(opts.DB))
		log = log.WithField("executionID", uuid.NewV4().String())
	}

	return log
}

func fileHK(opts FileOptions) logrus.Hook {
	fileHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel: newLogRotation(FileOptions{
			FolderPath: opts.FolderPath,
			FileName:   opts.FileName + ".info",
			MaxAge:     opts.MaxAge,
		}),
		logrus.ErrorLevel: newLogRotation(FileOptions{
			FolderPath: opts.FolderPath,
			FileName:   opts.FileName + ".error",
			MaxAge:     opts.MaxAge,
		}),
	})
	fileHook.SetFormatter(&logrus.JSONFormatter{})
	return fileHook
}

func newLogRotation(opts FileOptions) *logrotation.RotateLogs {
	pattern := "%Y-%m-%d"
	separator := "."

	path := opts.FolderPath
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return logrotation.New(
		path+opts.FileName+separator+pattern+separator+"log",
		logrotation.WithLinkName(path+opts.FileName+separator+"log"),
		logrotation.WithMaxAge(opts.MaxAge),
	)
}

func stackHK() logrus.Hook {
	return stackhook.NewHook(logrus.AllLevels, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
}

func mongoHK(opts MongoOptions) logrus.Hook {
	var mongoHook logrus.Hook
	if "" != opts.User && "" != opts.Password {
		var err error
		mongoHook, err = mgorus.NewHookerWithAuth(opts.Address+":"+opts.Port, opts.DB, opts.Collection, opts.User, opts.Password)
		if nil != err {
			panic(errors.Wrap(err, "Connecting to Mongo DB"))
		}
	} else {
		var err error
		mongoHook, err = mgorus.NewHooker(opts.Address+":"+opts.Port, opts.DB, opts.Collection)
		if nil != err {
			panic(errors.Wrap(err, "Connecting to Mongo DB"))
		}
	}
	return mongoHook
}
