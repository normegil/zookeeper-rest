package node

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	apiErrors "github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/normegil/zookeeper-rest/modules/tools"
	"github.com/pkg/errors"
)

const (
	CREATE_METHOD = "PUT"
	CREATE_PATH   = BASE_PATH
)

func (c Controller) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var body map[string]interface{}
	err := json.Unmarshal(tools.ToBytes(r.Body), &body)
	if nil != err {
		return errors.Wrap(err, "Error when unmarshalling Request body")
	}

	path, isString := body["path"].(string)
	if "" == path {
		return apiErrors.NewErrWithCode(40001, errors.New("Resource 'path' not specified."))
	}
	if !isString {
		return apiErrors.NewErrWithCode(40001, errors.New("Resource 'path' found but type not 'string'."))
	}

	content, isString := body["content"].(string)
	if !isString {
		return apiErrors.NewErrWithCode(40001, errors.New("Resource 'content' found but type not 'string'."))
	}

	c.Log().WithField("path", path).Debug("Creating node")
	return c.Zookeeper().Set(path, []byte(content))
}

const (
	UPDATE_METHOD = "POST"
	UPDATE_PATH   = BASE_PATH + "/:" + NODE_ID_PARAM_KEY
)

func (c Controller) update(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	path, err := c.path(p)
	if nil != err {
		return errors.Wrap(err, "Error when loading requested path")
	}

	var body map[string]interface{}
	err = json.Unmarshal(tools.ToBytes(r.Body), &body)
	if nil != err {
		return errors.Wrap(err, "Error when unmarshalling request body")
	}

	content, isString := body["content"].(string)
	if !isString {
		return apiErrors.NewErrWithCode(40001, errors.New("Resource 'content' found but type not 'string'."))
	}
	c.Log().WithField("path", path).Debug("Updating node")
	return c.Zookeeper().Set(path, []byte(content))
}
