package errors

import "github.com/pkg/errors"

type stacktracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

func stacks(err error) []errors.StackTrace {
	var stacktraces []errors.StackTrace
	errCauser, isCauser := err.(causer)
	if isCauser {
		stacktraces = stacks(errCauser.Cause())
	}

	stackErr, isStacktacer := err.(stacktracer)
	if isStacktacer {
		stacktraces = append(stacktraces, stackErr.StackTrace())
	}
	return stacktraces
}

func getErrWithCode(e error) ErrWithCode {
	found := SearchThroughCauses(e, func(e error) bool {
		_, isErrWithCode := e.(ErrWithCode)
		return isErrWithCode
	})

	if nil != found {
		return found.(ErrWithCode)
	}

	return NewErrWithCode(DEFAULT_CODE, e)
}

func SearchThroughCauses(e error, isSearched func(error) bool) error {
	if nil == e {
		return nil
	}

	isSearchedError := isSearched(e)
	if isSearchedError {
		return e
	}

	errWithCause, isErrWithCause := e.(causer)
	if !isErrWithCause || nil == errWithCause.Cause() {
		return nil
	}

	return SearchThroughCauses(errWithCause.Cause(), isSearched)
}
