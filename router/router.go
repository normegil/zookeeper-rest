package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

type Router struct {
	log    *logrus.Entry
	router *httprouter.Router
}

func New(log *logrus.Entry) *Router {
	return &Router{
		log:    log,
		router: httprouter.New(),
	}
}

func (r *Router) Register(routes []Route) error {
	for _, route := range routes {
		switch route.Method() {
		case "HEAD":
			r.router.HEAD(route.Path(), route.Handler())
		case "GET":
			r.router.GET(route.Path(), route.Handler())
		case "POST":
			r.router.POST(route.Path(), route.Handler())
		case "PUT":
			r.router.PUT(route.Path(), route.Handler())
		case "DELETE":
			r.router.DELETE(route.Path(), route.Handler())
		case "OPTIONS":
			r.router.OPTIONS(route.Path(), route.Handler())
		case "PATCH":
			r.router.PATCH(route.Path(), route.Handler())
		default:
			return errors.New("HTTP Method not supported {method: " + route.Method() + "; path: " + route.Path() + "}")
		}
	}
	return nil
}

func (r *Router) Listen(port int) error {
	r.log.WithField("port", port).Info("Launching server")
	if err := http.ListenAndServe(":"+strconv.Itoa(port), r.router); nil != err {
		return err
	}
	return nil
}
