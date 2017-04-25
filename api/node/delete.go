package node

import (
	"errors"
	"net/http"

	errPkg "github.com/pkg/errors"

	"github.com/julienschmidt/httprouter"
)

func (c Controller) remove(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	path, err := c.path(params)
	if nil != err {
		return errPkg.Wrap(err, "Could not get node's path: ")
	} else if "" == path {
		return errors.New("No path corresponding to ID")
	}

	c.Log().WithField("path", path).Debug("Deleting node")
	err = c.Zookeeper().Delete(path)
	if nil != err {
		return errPkg.Wrapf(err, "Could not remove %s: ", path)
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
