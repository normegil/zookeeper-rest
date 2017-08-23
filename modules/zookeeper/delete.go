package zookeeper

import "github.com/pkg/errors"

func (z Zookeeper) Delete(path string, recursive bool) error {
	client, err := z.client()
	if nil != err {
		return errors.Wrap(err, "Obtaining Zookeeper client")
	}
	defer client.Close()

	return deleteNode(client, path, recursive)
}

func deleteNode(client ZookeeperClient, path string, recursive bool) error {
	_, stat, err := client.Get(path)
	if nil != err {
		return errors.Wrapf(err, "Loading statistics of %s", path)
	}

	if recursive {
		childs, _, err := client.Children(path)
		if err != nil {
			return errors.Wrapf(err, "Getting childs od %s", path)
		}
		for _, child := range childs {
			fullChildPath := path + "/" + child
			err := deleteNode(client, fullChildPath, recursive)
			if err != nil {
				return errors.Wrapf(err, "Deleting %s", fullChildPath)
			}
		}
	}

	if err = client.Delete(path, stat.Version()); err != nil {
		return errors.Wrapf(err, "Deleting %s", path)
	}
	return nil
}
