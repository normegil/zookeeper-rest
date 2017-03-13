package router

import "github.com/julienschmidt/httprouter"

type Controller interface {
	Routes() []Route
}

type Route interface {
	Method() string
	Path() string
	Handler() httprouter.Handle
}

type HttpRoute struct {
	method  string
	path    string
	handler httprouter.Handle
}

func (r HttpRoute) Method() string {
	return r.method
}

func (r HttpRoute) Path() string {
	return r.path
}

func (r HttpRoute) Handler() httprouter.Handle {
	return r.handler
}

func NewRoute(method, path string, handler httprouter.Handle) *HttpRoute {
	return &HttpRoute{
		method:  method,
		path:    path,
		handler: handler,
	}
}
