package zookeeper

import (
	"encoding/json"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const ZOOKEEPER_CONFIG_ROOT_NODE = "/zk-rest"
const ZOOKEEPER_NODE_PATH_STORAGE string = ZOOKEEPER_CONFIG_ROOT_NODE + "/pathIDs"

func (z Zookeeper) Path(id string) (string, error) {
	if "" == id {
		z.Log().Debug("No ID: Returning root")
		return "/", nil
	}

	err := z.refreshPathList()
	if nil != err {
		return "", errors.Wrap(err, "Refreshing path-id associations")
	}

	content, err := z.Content(ZOOKEEPER_NODE_PATH_STORAGE)
	if nil != err {
		return "", errors.Wrap(err, "Get content of paths cache")
	}

	paths := make(map[string]string)
	err = json.Unmarshal([]byte(content), &paths)
	if nil != err {
		return "", errors.Wrap(err, "Unmarshalling path cache")
	}

	path, ok := paths[id]
	if !ok {
		return "", errors.New("Path not found for ID: " + id)
	}

	z.Log().WithField("id", id).WithField("path", path).Debug("ZK: Path found")
	return path, nil
}

func (z Zookeeper) refreshPathList() error {
	exist, err := z.Exist(ZOOKEEPER_CONFIG_ROOT_NODE)
	if nil != err {
		return errors.Wrap(err, "Check existence of root node for application configuration")
	}
	if !exist {
		z.Log().WithField("configPath", ZOOKEEPER_CONFIG_ROOT_NODE).Info("ZK: Root config node doesn't exist")
		err := z.Set(ZOOKEEPER_CONFIG_ROOT_NODE, []byte{})
		if nil != err {
			return errors.Wrap(err, "Creating root configuration node")
		}
	}

	exist, err = z.Exist(ZOOKEEPER_NODE_PATH_STORAGE)
	if nil != err {
		return errors.Wrap(err, "Check existence of path-id association node")
	}
	paths := make(map[string]string)
	if !exist {
		z.Log().WithField("configPath", ZOOKEEPER_NODE_PATH_STORAGE).Info("ZK: Paths config node doesn't exist")
		err = z.Set(ZOOKEEPER_NODE_PATH_STORAGE, []byte{})
		if nil != err {
			return errors.Wrap(err, "Create path-id association node")
		}
	} else {
		content, err := z.Content(ZOOKEEPER_NODE_PATH_STORAGE)
		if nil != err {
			return errors.Wrap(err, "Get content of path-id association node")
		}

		if "" != content {
			err = json.Unmarshal([]byte(content), &paths)
			if nil != err {
				return errors.Wrap(err, "Unmarshall content of path-id association node")
			}
		}
	}

	z.Log().Info("Retreiving all existing paths in ZK root")
	subPaths, err := z.loadSubPath("/")
	if nil != err {
		return errors.Wrap(err, "Browse through Zookeeper tree")
	}
	subPaths = append(subPaths, "/")
	z.Log().Info("Tree paths retrieved")

	for _, subPath := range subPaths {
		if !isRegistered(subPath, paths) {
			paths[uuid.NewV4().String()] = subPath
		}
	}

	toSave, err := json.Marshal(paths)
	if nil != err {
		return errors.Wrap(err, "Marshall new content for path-id association node")
	}

	z.Log().WithField("configNode", ZOOKEEPER_NODE_PATH_STORAGE).Info("Saving extra paths into paths config node")
	return errors.Wrap(z.Set(ZOOKEEPER_NODE_PATH_STORAGE, toSave), "Saving path-id association node with new contents")
}

func (z Zookeeper) loadSubPath(path string) ([]string, error) {
	childs, err := z.Children(path)
	if nil != err {
		return nil, err
	}

	toReturn := make([]string, len(childs))
	for _, child := range childs {
		childFullPath := path + "/" + child
		if "/" == path {
			childFullPath = "/" + child
		}
		toReturn = append(toReturn, childFullPath)
		subChilds, err := z.loadSubPath(childFullPath)
		if nil != err {
			return nil, err
		}
		toReturn = append(toReturn, subChilds...)
	}
	return toReturn, nil
}

func isRegistered(path string, registeredPaths map[string]string) bool {
	for _, registeredPath := range registeredPaths {
		if registeredPath == path {
			return true
		}
	}
	return false
}

func (z Zookeeper) ID(path string) (string, error) {
	if "" == path {
		z.Log().Debug("No Path: using root")
		path = "/"
	}

	paths, err := z.IDs([]string{path})
	if nil != err {
		return "", err
	}
	toReturn, found := paths[path]
	if !found {
		return "", errors.New("ID not found for path " + path)
	}
	return toReturn, nil
}

func (z Zookeeper) IDs(paths []string) (map[string]string, error) {
	z.Log().WithField("paths", paths).Debug("Research IDs for paths")
	err := z.refreshPathList()
	if nil != err {
		return nil, errors.Wrap(err, "Refreshing path-id associations")
	}

	content, err := z.Content(ZOOKEEPER_NODE_PATH_STORAGE)
	if nil != err {
		return nil, errors.Wrap(err, "Load content of paths-id association")
	}

	var associations map[string]string
	err = json.Unmarshal([]byte(content), &associations)
	if nil != err {
		return nil, errors.Wrap(err, "Unmarshalling path cache")
	}

	ids := make(map[string]string)
	for _, path := range paths {
		for key, loadedPath := range associations {
			if path == loadedPath {
				ids[path] = key
				break
			}
		}
	}

	return ids, nil
}
