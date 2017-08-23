package node

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const NODE_ID_PARAM_KEY = "nodeID"

func (c Controller) path(params httprouter.Params) (string, error) {
	path := params.ByName(NODE_ID_PARAM_KEY)

	if "" == path {
		return "/", nil
	}
	id, err := uuid.FromString(path)
	if nil != err {
		return "", errors.Wrapf(err, "Parse parameter %s into UUID", NODE_ID_PARAM_KEY)
	}

	p, err := c.Zookeeper().Path(id.String())
	if nil != err {
		return "", errors.Wrap(err, "Retrieving path from Zookeeper")
	}
	c.Log().WithField("id", id).WithField("path", p).Debug("Loaded path from UUID")
	return p, nil
}
