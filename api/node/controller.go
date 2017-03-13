package node

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
		router.NewRoute("GET", "/rest/node", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			fmt.Fprintf(w, "%+v", r.URL.Query())
		}),
		router.NewRoute("GET", "/rest/node/childs", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			fmt.Fprintf(w, "%+v", r.URL.Query())
		}),
	}
}
