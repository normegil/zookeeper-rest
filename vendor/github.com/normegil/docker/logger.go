package docker

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type DefaultLogger struct{}

func (l DefaultLogger) Print(v ...interface{})                {}
func (l DefaultLogger) Printf(fomat string, v ...interface{}) {}
