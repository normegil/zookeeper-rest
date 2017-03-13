package node

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/normegil/zookeeper-rest/modules/environment"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
	"github.com/normegil/zookeeper-rest/router"
)

type Controller struct {
	environment.Env
}

const ADDRESS string = "127.0.0.1"

func (c Controller) Routes() []router.Route {
	return []router.Route{
		router.NewRoute("GET", "/rest/node", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			children, err := zookeeper.Zookeeper{ADDRESS}.Children(r.URL.Query().Get("path"))
			if nil != err {
				c.Log().WithError(err).Error("Error detected")
				fmt.Fprintf(w, "Error")
				return
			}
			fmt.Fprintf(w, "%+v", children)
		}),
		router.NewRoute("GET", "/rest/node/childs", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			fmt.Fprintf(w, "%+v", r.URL.Query())
		}),
	}
}
