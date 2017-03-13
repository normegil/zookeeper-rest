package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/normegil/zookeeper-rest/modules/environment"
	"github.com/normegil/zookeeper-rest/router"
)

type Controller struct {
	environment.Env
}

func (c Controller) Routes() []router.Route {
	return []router.Route{
		router.NewRoute("GET", "/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			fmt.Fprintf(w, "Test")
		}),
	}
}
