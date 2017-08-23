package router

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/normegil/zookeeper-rest/modules/database/mongo"
	"github.com/normegil/zookeeper-rest/modules/environment"
	"github.com/normegil/zookeeper-rest/modules/middleware"
	"github.com/pkg/errors"
)

type Router struct {
	environment.Env
	router *httprouter.Router
}

func New(env environment.Env) *Router {
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
	dao := &mongo.MongoUserDAO{r.Env.Session(), ""}
	handler := middleware.URLContructor(middleware.RequestLogger(r.Env.Log(), middleware.RequestAuthenticator(r.Env.Log(), dao, r.router)))

	r.Log().WithField("port", port).Info("Launching server")
	if err := http.ListenAndServe(":"+strconv.Itoa(port), handler); nil != err {
		return errors.Wrapf(err, "Error while Listening on %d", port)
	}
	return nil
}
