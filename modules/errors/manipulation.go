package errors

import "github.com/pkg/errors"

type stacktracer interface {
	StackTrace() errors.StackTrace
}

type causer interface{
	Cause() error
}

func stacks(err error) []errors.StackTrace {
	var stacktraces []errors.StackTrace
	errCauser, isCauser := err.(causer)
	if isCauser {
		stacktraces = stacks(errCauser.Cause())
	}

	stackErr, isStacktacer := err.(stacktracer)
	if (isStacktacer) {
		stacktraces = append(stacktraces, stackErr.StackTrace())
	}
	return stacktraces
}
