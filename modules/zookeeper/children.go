package zookeeper

import "github.com/pkg/errors"

func (z Zookeeper) Children(path string) ([]string, error) {
	client, err := z.client()
	if nil != err {
		return nil, errors.Wrap(err, "Could not get Zookeeper Client")
	}
	defer client.Close()

	z.Log().WithField("parentPath", path).Debug("Load childs")
	children, _, err := client.Children(path)
	if nil != err {
		return nil, errors.Wrapf(err, "Loading childs for %s", path)
	}
	return children, nil
}
