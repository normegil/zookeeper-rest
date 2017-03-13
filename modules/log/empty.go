package log

type VoidLogger struct{}

func (v VoidLogger) Printf(string, ...interface{}) {}
