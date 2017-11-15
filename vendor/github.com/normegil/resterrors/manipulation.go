package resterrors

import "github.com/pkg/errors"

type Stacktracer interface {
	StackTrace() errors.StackTrace
}

type Causer interface {
	Cause() error
}

func Stacks(err error) []errors.StackTrace {
	var stacktraces []errors.StackTrace
	errCauser, isCauser := err.(Causer)
	if isCauser {
		stacktraces = Stacks(errCauser.Cause())
	}

	stackErr, isStacktacer := err.(Stacktracer)
	if isStacktacer {
		stacktraces = append(stacktraces, stackErr.StackTrace())
	}
	return stacktraces
}

func getErrWithCode(e error, defaultCode int) ErrWithCode {
	found := SearchThroughCauses(e, func(e error) bool {
		_, isErrWithCode := e.(ErrWithCode)
		return isErrWithCode
	})

	if nil != found {
		return found.(ErrWithCode)
	}

	return NewErrWithCode(defaultCode, e)
}

func SearchThroughCauses(e error, isSearched func(error) bool) error {
	if nil == e {
		return nil
	}

	isSearchedError := isSearched(e)
	if isSearchedError {
		return e
	}

	errWithCause, isErrWithCause := e.(Causer)
	if !isErrWithCause || nil == errWithCause.Cause() {
		return nil
	}

	return SearchThroughCauses(errWithCause.Cause(), isSearched)
}
