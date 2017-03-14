package node

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
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
	}
}

func (c Controller) loadContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	content, err := c.Zookeeper().Content(r.URL.Query().Get("path"))
	if nil != err {
		c.handleError(w, err)
		return
	}

	response, err := json.Marshal(content)
	if nil != err {
		c.handleError(w, err)
		return
	}
	fmt.Fprintf(w, string(response))
}

func (c Controller) loadChildren(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	children, err := c.Zookeeper().Children(r.URL.Query().Get("path"))
	if nil != err {
		c.handleError(w, err)
		return
	}

	response, err := json.Marshal(children)
	if nil != err {
		c.handleError(w, err)
		return
	}

	fmt.Fprintf(w, string(response))
}

func (c Controller) handleError(w http.ResponseWriter, err error) {
	newErr := err
	switch err {
	case zk.ErrNoServer:
		newErr = errors.NewErrWithCode(50301, err)
	}
	errors.Handler{c.Log()}.Handle(w, newErr)
}
