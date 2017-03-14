package node

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
