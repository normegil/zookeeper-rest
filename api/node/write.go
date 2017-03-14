package node

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	apiErrors "github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/normegil/zookeeper-rest/modules/tools"
	"github.com/samuel/go-zookeeper/zk"
)

type node struct {
	Path    string
	Data    string
	Version int
}

func (c Controller) createOrUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	path := r.URL.Query().Get("path")
	if "" == path {
		c.create(w, r, p)
		return
	}
	c.update(w, r, p)
}

func (c Controller) create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var toCreate node
	err := json.Unmarshal(tools.ToBytes(r.Body), &toCreate)
	if nil != err {
		c.handleError(w, err)
		return
	}

	if "" == toCreate.Path {
		err = apiErrors.NewErrWithCode(40001, errors.New("Resource path not specified."))
		c.handleError(w, err)
		return
	}

	acls := zk.WorldACL(zk.PermAll)
	if err := c.Zookeeper().Create(toCreate.Path, []byte(toCreate.Data), acls); nil != err {
		c.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c Controller) update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var toUpdate node
	err := json.Unmarshal(tools.ToBytes(r.Body), &toUpdate)
	if nil != err {
		c.handleError(w, err)
		return
	}

	toUpdate.Path = r.URL.Query().Get("path")
	if "" == toUpdate.Path {
		err = apiErrors.NewErrWithCode(40001, errors.New("Resource path not specified."))
		c.handleError(w, err)
		return
	}

	err = c.Zookeeper().Set(toUpdate.Path, []byte(toUpdate.Data), toUpdate.Version)
	if nil != err {
		c.handleError(w, err)
		return
	}
}
