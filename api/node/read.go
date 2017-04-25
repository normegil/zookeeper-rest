package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	jsonUtils "github.com/normegil/zookeeper-rest/modules/json"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type nodeResponse struct {
	ID      string
	URL     jsonUtils.JSONURL
	Path    string
	Content string
	Childs  map[string]jsonUtils.JSONURL
}

func (c Controller) load(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	paramID := params.ByName(KeyNodeID)
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

	node, err := c.Zookeeper().Load(path)
	if nil != err {
		return errors.Wrap(err, "Loading node content")
	}

	c.Log().WithField("Path", path).Debug("Node read")

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

	ids, err := c.Zookeeper().IDs(node.Childs)
	if nil != err {
		return errors.Wrap(err, "Loading childs paths IDs")
	}

	c.Log().WithField("childs", ids).Debug("Child IDs loaded")
	childURLs := make(map[string]jsonUtils.JSONURL)
	for path, key := range ids {
		newURL, err := url.Parse(baseURL + "/" + key)
		childURLs[path] = jsonUtils.JSONURL(*newURL)
		if nil != err {
			return errors.Wrapf(err, "Cannot parse URL from %s", path)
		}
	}

	nodeResp := nodeResponse{
		ID:      id,
		URL:     jsonUtils.JSONURL(*currentURL),
		Path:    node.Path,
		Content: node.Content,
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
