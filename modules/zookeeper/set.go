package zookeeper

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

func (z Zookeeper) Set(path string, content []byte) error {
	client, err := z.client()
	if nil != err {
		return errors.Wrap(err, "Connecting to zookeeper")
	}
	defer client.Close()
	//client.SetLogger(log.VoidLogger{})

	exist, _, err := client.Exists(path)
	if nil != err {
		return errors.Wrap(err, "Testing path for existence")
	}

	z.Log().WithField("path", path).WithField("Creation", !exist).Debug("Node Creation/Update")
	if exist {
		_, stat, err := client.Get(path)
		if nil != err {
			return errors.Wrap(err, "Getting path version")
		}
		_, err = client.Set(path, content, stat.Version())
		if err != nil {
			return errors.Wrapf(err, "Setting %s content", path)
		}
	} else {
		var noFlag int32 = 0
		if err := recursiveCreate(client, path, content, noFlag, zk.WorldACL(zk.PermAll)); err != nil {
			return errors.Wrapf(err, "Creating path %s", path)
		}
	}
	return nil
}

func recursiveCreate(client ZookeeperClient, path string, content []byte, flags int32, acls []zk.ACL) error {
	exist, _, err := client.Exists(path)
	if nil != err {
		return errors.Wrap(err, "Testing path for existence")
	}

	if exist {
		return nil
	}

	if "/" != path {
		parentPath := parent(path)
		if err := recursiveCreate(client, parentPath, []byte(""), flags, acls); err != nil {
			return errors.Wrapf(err, "Recursive creating parent %s", parentPath)
		}
	}

	_, err = client.Create(path, content, flags, acls)
	if err != nil {
		return errors.Wrapf(err, "Creating path %s", path)
	}
	return nil
}

func parent(path string) string {
	splittedPath := strings.Split(path, "/")
	if 0 == len(splittedPath) || 1 == len(splittedPath) || 2 == len(splittedPath) {
		return "/"
	}

	parentSplittedPath := splittedPath[:len(splittedPath)-1]
	return strings.Join(parentSplittedPath, "/")
}
