package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/normegil/zookeeper-rest/model"
	"github.com/normegil/zookeeper-rest/router/middleware"
)

type Router struct {
	model.Env
	router *httprouter.Router
}

func New(env model.Env) *Router {
	return &Router{
		Env:    env,
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
	handler := middleware.RequestLogger(r.Env, r.router)

	r.Log().WithField("port", port).Info("Launching server")
	if err := http.ListenAndServe(":"+strconv.Itoa(port), handler); nil != err {
		return err
	}
	return nil
}
