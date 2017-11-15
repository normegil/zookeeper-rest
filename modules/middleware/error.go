package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/normegil/resterrors"
	definitions "github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/sirupsen/logrus"
)

type HttpHandlerFunc func(http.ResponseWriter, *http.Request, httprouter.Params) error
type ErrorHandlerFunc func(http.ResponseWriter, error)
type ErrorCodeAssignerFunc func(err error) error

type ErrorHandler struct {
	handler         HttpHandlerFunc
	errCodeAssigner ErrorCodeAssignerFunc
	errHandler      ErrorHandlerFunc

	logger logrus.FieldLogger
}

func (e ErrorHandler) Handle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := e.handler(w, r, p)
	if nil == err {
		return
	}

	if nil != e.errCodeAssigner {
		err = e.errCodeAssigner(err)
	}

	if nil != e.errHandler {
		e.errHandler(w, err)
		return
	}

	stack := resterrors.Stacks(err)
	e.logger.WithField("error", err).WithField("stack", stack).Info("Error")
	resterrors.Handler{definitions.Definitions(), definitions.DEFAULT_ERROR_CODE}.Handle(w, err)
}

type ErrorHandlerFactory struct {
	ErrorHandlerFunc
	ErrorCodeAssignerFunc
	Logger logrus.FieldLogger
}

func (f ErrorHandlerFactory) New(handler HttpHandlerFunc) *ErrorHandler {
	h := ErrorHandler{
		handler:         handler,
		errCodeAssigner: f.ErrorCodeAssignerFunc,
		errHandler:      f.ErrorHandlerFunc,
		logger:          f.Logger,
	}
	return &h
}
