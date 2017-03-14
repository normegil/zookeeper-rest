package node

import (
	"net/http"

	"github.com/normegil/zookeeper-rest/modules/environment"
	"github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/normegil/zookeeper-rest/router"
	"github.com/samuel/go-zookeeper/zk"
)

type Controller struct {
	environment.Env
}

const ADDRESS string = "127.0.0.1"

func (c Controller) Routes() []router.Route {
	return []router.Route{
		router.NewRoute("GET", "/rest/node/childs", c.loadChildren),
		router.NewRoute("GET", "/rest/node", c.loadContent),
		router.NewRoute("POST", "/rest/node", c.createOrUpdate),
		router.NewRoute("PUT", "/rest/node", c.createOrUpdate),
	}
}

func (c Controller) handleError(w http.ResponseWriter, err error) {
	newErr := err
	switch err {
	case zk.ErrNoServer:
		newErr = errors.NewErrWithCode(50301, err)
	}
	errors.Handler{c.Log()}.Handle(w, newErr)
}
