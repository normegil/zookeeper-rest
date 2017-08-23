package node

import (
	"net/http"
	"strconv"

	errPkg "github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"

	"github.com/julienschmidt/httprouter"
)

const (
	DELETE_METHOD = "DELETE"
	DELETE_PATH   = BASE_PATH + "/:" + NODE_ID_PARAM_KEY
)

func (c Controller) remove(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	path, err := c.path(params)
	if nil != err {
		return errors.Wrap(err, "Could not get node's path: ")
	} else if "" == path {
		return errors.New("No path corresponding to ID")
	} else if "/" == path {
		return errPkg.NewErrWithCode(40004, errors.New("Cannot remove root path"))
	}

	recursiveStr := r.URL.Query().Get("recursive")
	recursive := false
	if "" != recursiveStr {
		recursive, err = strconv.ParseBool(recursiveStr)
		if err != nil {
			return errPkg.NewErrWithCode(40003, errors.Wrapf(err, "'recursive' query parameter could not be parsed (Wrong value %s)", recursiveStr))
		}
	}

	c.Log().WithField("path", path).WithField("recursive", recursive).Info("Deleting node")
	if err = c.Zookeeper().Delete(path, recursive); nil != err {
		toReturn := errors.Wrapf(err, "Could not remove %s: ", path)
		found := errPkg.SearchThroughCauses(err, func(e error) bool {
			return zk.ErrNotEmpty == e
		})
		if nil != found {
			return errPkg.NewErrWithCode(40005, toReturn)
		}
		return toReturn
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
