package log

import (
	"strings"
	"time"

	stackhook "github.com/Gurpartap/logrus-stack"
	logrotation "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/weekface/mgorus"
)

type Options struct {
	// Verbose controls the verbosity of the log (Set to true if you want a DebugLevel logger)
	Verbose bool
	File    FileOptions
	DB      MongoOptions
}

// FileOptions controls the options that revolve around file logging. If it's not specified, file won't be created
type FileOptions struct {
	//FolderPath is the path to the folder where the log files will be stored
	FolderPath string
	//FileName is the filename prefix that will be used to create the log file. Suffix will be added to differentiate between log file types.
	FileName string
	//MaxAge is the dureation between rotations of the log files
	MaxAge time.Duration
}

// MongoOptions controls the options that revolve around file logging. If it's not specified, file won't be created
type MongoOptions struct {
	URL      string
	Database string
	User     string
	Password string
}

// New create a new logrus logger & configure it based on Options
func New(opts Options) (*logrus.Entry, error) {
	log := logrus.NewEntry(logrus.New())
	if opts.Verbose {
		log.Logger.Level = logrus.DebugLevel
	}

	log.Logger.Hooks.Add(stackHK())

	if "" != opts.File.FileName {
		hook, err := fileHK(opts.File)
		if err != nil {
			return nil, err
		}
		log.Logger.Hooks.Add(hook)
	}

	if "" != opts.DB.URL {
		log.Logger.Hooks.Add(mongoHK(opts.DB))
		log = log.WithField("executionID", uuid.NewV4().String())
	}

	return log, nil
}

func fileHK(opts FileOptions) (logrus.Hook, error) {
	infoRotation, err := newLogRotation(FileOptions{
		FolderPath: opts.FolderPath,
		FileName:   opts.FileName + ".info",
		MaxAge:     opts.MaxAge,
	})
	if err != nil {
		return nil, err
	}

	errorRotation, err := newLogRotation(FileOptions{
		FolderPath: opts.FolderPath,
		FileName:   opts.FileName + ".error",
		MaxAge:     opts.MaxAge,
	})
	if err != nil {
		return nil, err
	}

	fileHook := lfshook.NewHook(lfshook.WriterMap{logrus.InfoLevel: infoRotation, logrus.ErrorLevel: errorRotation})
	fileHook.SetFormatter(&logrus.JSONFormatter{})
	return fileHook, nil
}

func newLogRotation(opts FileOptions) (*logrotation.RotateLogs, error) {
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
	collection := "log"
	if "" != opts.User && "" != opts.Password {
		var err error
		mongoHook, err = mgorus.NewHookerWithAuth(opts.URL, opts.Database, collection, opts.User, opts.Password)
		if nil != err {
			panic(errors.Wrap(err, "Connecting to Mongo DB"))
		}
	} else {
		var err error
		mongoHook, err = mgorus.NewHooker(opts.URL, opts.Database, collection)
		if nil != err {
			panic(errors.Wrap(err, "Connecting to Mongo DB"))
		}
	}
	return mongoHook
}
