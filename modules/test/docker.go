package test

import (
	"testing"

	"github.com/normegil/docker"
)

type TestLogger struct {
	t testing.TB
}

func (l *TestLogger) Print(v ...interface{}) {
	l.t.Log(v...)
}

func (l *TestLogger) Printf(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}

func NewDocker(t testing.TB, options docker.Options) (*docker.ContainerInfo, func()) {
	if nil == options.Logger {
		options.Logger = &TestLogger{t}
	}
	info, closeFn, err := docker.New(options)
	if err != nil {
		t.Fatal(err)
	}
	return info, func() {
		err := closeFn()
		if nil != err {
			t.Fatal(err)
		}
	}
}
