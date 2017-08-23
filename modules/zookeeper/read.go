package zookeeper

import "github.com/pkg/errors"

func (z Zookeeper) Content(path string) ([]byte, error) {
	client, err := z.client()
	if nil != err {
		return nil, err
	}
	defer client.Close()

	content, _, err := client.Get(path)
	if nil != err {
		return nil, errors.Wrapf(err, "Loading node %s", path)
	}
	return content, nil
}

func (z Zookeeper) Exist(path string) (bool, error) {
	client, err := z.client()
	if nil != err {
		return false, err
	}
	defer client.Close()

	exist, _, err := client.Exists(path)
	if nil != err {
		return false, errors.Wrapf(err, "Checking existence of %s", path)
	}
	return exist, nil
}

func (z Zookeeper) Load(path string, recursive bool) (*nodeImpl, error) {
	if recursive {
		return NewRecursiveNode(path, z)
	}

	client, err := z.client()
	if nil != err {
		return nil, errors.Wrap(err, "Connecting to Zookeeper")
	}
	defer client.Close()

	content, _, err := client.Get(path)
	if nil != err {
		return nil, errors.Wrapf(err, "Loading node %s", path)
	}

	childs, _, err := client.Children(path)
	if nil != err {
		return nil, errors.Wrapf(err, "Loading childs of %s", path)
	}

	return NewNode(path, content, childs), nil
}
