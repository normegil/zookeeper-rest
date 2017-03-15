package node

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (c Controller) remove(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := c.Zookeeper().Delete(r.URL.Query().Get("path"))
	if nil != err {
		c.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
