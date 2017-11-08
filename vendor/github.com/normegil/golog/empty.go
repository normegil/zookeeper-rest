package log

type SimpleLogger interface {
	Printf(string, ...interface{})
}

type VoidLogger struct{}

func (v VoidLogger) Printf(string, ...interface{}) {}
