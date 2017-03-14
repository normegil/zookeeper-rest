package node

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	apiErrors "github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/normegil/zookeeper-rest/modules/tools"
)

func (c Controller) createOrUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	path := r.URL.Query().Get("path")
	if "" == path {
		err := apiErrors.NewErrWithCode(40001, errors.New("Resource path not specified."))
		c.handleError(w, err)
		return
	}

	if err := c.Zookeeper().Set(path, tools.ToBytes(r.Body)); nil != err {
		c.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
