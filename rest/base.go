package rest

import (
	"net/http"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/normegil/zookeeper-rest/router"
)

type Controller struct {
	Log *logrus.Entry
}

func (c Controller) Routes() []router.Route {
	return []router.Route{
		router.NewRoute("GET", "/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			c.Log.Info("Request for root")
			fmt.Fprintf(w, "Test")
		}),
	}
}
