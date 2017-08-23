package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/normegil/zookeeper-rest/modules/formats"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	GET_METHOD     = "GET"
	GET_PATH       = BASE_PATH + "/:" + NODE_ID_PARAM_KEY
	GET_ALL_METHOD = GET_METHOD
	GET_ALL_PATH   = BASE_PATH
)

func (c Controller) load(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	paramID := params.ByName(NODE_ID_PARAM_KEY)
	path := "/"
	if "" != paramID {
		id, err := uuid.FromString(paramID)
		if nil != err {
			return errors.Wrap(err, "Cannot generate UUID")
		}
		path, err = c.Zookeeper().Path(id.String())
		if nil != err {
			return errors.Wrap(err, "Cannot get path from Request parameter")
		}
	}

	node, err := c.Zookeeper().Load(path, false)
	if nil != err {
		return errors.Wrap(err, "Loading node content")
	}
	c.Log().WithField("node", node).Debugf("Node loaded: %s", path)

	baseURL := "http://" + r.Host + r.URL.Path
	if "" != paramID {
		splitted := strings.Split(baseURL, "/"+paramID)
		c.Log().WithField("splitted", splitted).WithField("sep", "/"+paramID).WithField("base", baseURL).Debug("Split URL")
		baseURL = splitted[0]
	}

	id, err := c.Zookeeper().ID(path)
	if nil != err {
		return errors.Wrap(err, "Loading current path ID")
	}
	urlAsStr := baseURL + "/" + id
	currentURL, err := url.Parse(urlAsStr)
	if nil != err {
		return errors.Wrapf(err, "Parsing current URL (%s)", urlAsStr)
	}

	childPaths := make([]string, len(node.Childs()))
	childNodes := node.Childs()
	c.Log().WithField("ChildNodes", childNodes).Debugf("ChildNodes loaded: %s", path)
	for i, child := range childNodes {
		childPaths[i] = child.Path()
	}

	ids, err := c.Zookeeper().IDs(childPaths)
	if nil != err {
		return errors.Wrap(err, "Loading childs paths IDs")
	}
	c.Log().WithField("childs", ids).WithField("paths", childPaths).Debugf("Child IDs: %s", path)

	childURLs := make(map[string]formats.URL)
	for path, key := range ids {
		newURL, err := url.Parse(baseURL + "/" + key)
		childURLs[path] = formats.URL{newURL}
		if nil != err {
			return errors.Wrapf(err, "Cannot parse URL from %s", path)
		}
	}

	nodeResp := nodeResponse{
		ID:      id,
		URL:     formats.URL{currentURL},
		Path:    node.Path(),
		Content: string(node.Content()),
		Childs:  childURLs,
	}

	response, err := json.Marshal(nodeResp)
	if nil != err {
		return errors.Wrap(err, "Marshalling response")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(response))
	return nil
}

type nodeResponse struct {
	ID      string
	URL     formats.URL
	Path    string
	Content string
	Childs  map[string]formats.URL
}

func (n nodeResponse) Equals(other nodeResponse) bool {
	if n.ID != other.ID {
		return false
	}
	if n.URL != other.URL {
		return false
	}
	if n.Path != other.Path {
		return false
	}
	if n.Content != other.Content {
		return false
	}
	if len(n.Childs) != len(other.Childs) {
		return false
	}
	for key, child := range n.Childs {
		if child != other.Childs[key] {
			return false
		}
	}
	return true
}
