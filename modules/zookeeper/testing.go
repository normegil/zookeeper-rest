package zookeeper

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	log "github.com/normegil/golog"
	"github.com/normegil/zookeeper-rest/modules/test"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
)

func Test_NewZookeeperClient(t testing.TB) (Zookeeper, func()) {
	if testing.Short() {
		zk := Zookeeper{
			Logger:             logrus.New(),
			SubstitutionClient: Test_NewInMemoryClient(),
		}
		err := zk.Set("/zookeeper/quota", []byte("")) //Real zookeeper instance always have the /zookeeper/quota node
		if err != nil {
			t.Fatal(errors.Wrapf(err, "Creating /zookeeper/quota in memory"))
		}
		return zk, func() {}
	} else {
		contactInfo, closeFn := test.NewZookeeper(t)
		return Zookeeper{
			Logger:  logrus.New(),
			Address: contactInfo,
		}, closeFn
	}
}

func Test_CleanInsert(t testing.TB, z Zookeeper, node Node) {
	if err := Test_CleanUp(z); err != nil {
		t.Fatal(errors.Wrapf(err, "Cleaning zookeeper"))
	}
	if err := Test_InsertNode(z, node); err != nil {
		t.Fatal(errors.Wrapf(err, "Inserting Nodes into Zookeeper"))
	}
}

func Test_InsertNode(z Zookeeper, node Node) error {
	err := z.Set(node.Path(), node.Content())
	if nil != err {
		return errors.Wrap(err, "Creating node "+node.Path())
	}

	identified, ok := node.(Identified)
	if ok && "" != identified.ID() {
		if err := Test_RegisterNodeID(z, identified.ID(), node.Path()); err != nil {
			return errors.Wrapf(err, "Registering associactions (%s;%s)", identified.ID(), node.Path())
		}
	}

	for _, child := range node.Childs() {
		if err := Test_InsertNode(z, child); err != nil {
			return errors.Wrapf(err, "Inserting node %s", node.Path)
		}
	}
	return nil
}

func Test_RegisterNodeID(z Zookeeper, id string, path string) error {
	exist, err := z.Exist(ZOOKEEPER_NODE_PATH_STORAGE)
	if err != nil {
		return errors.Wrapf(err, "Testing existence %s", ZOOKEEPER_NODE_PATH_STORAGE)
	}
	pathAssociations := make(map[string]string)
	if exist {
		node, err := z.Load(ZOOKEEPER_NODE_PATH_STORAGE, false)
		if nil != err {
			return errors.Wrapf(err, "Getting path associations")
		}
		if err := json.Unmarshal(node.Content(), &pathAssociations); err != nil {
			return errors.Wrapf(err, "Unarshal content of %s", ZOOKEEPER_NODE_PATH_STORAGE)
		}
	}
	for id, pathAssociated := range pathAssociations {
		if pathAssociated == path {
			delete(pathAssociations, id)
		}
	}
	pathAssociations[id] = path
	content, err := json.Marshal(pathAssociations)
	if err != nil {
		return errors.Wrapf(err, "Marshal %s content (%+v)", ZOOKEEPER_NODE_PATH_STORAGE, pathAssociations)
	}
	if err := z.Set(ZOOKEEPER_NODE_PATH_STORAGE, content); err != nil {
		return errors.Wrapf(err, "Setting content of %s", ZOOKEEPER_NODE_PATH_STORAGE)
	}
	return nil
}

func Test_CleanUp(z Zookeeper) error {
	childs, err := z.Children("/")
	for _, child := range childs {
		if "zookeeper" == child { // '/zookeeper' connot be deleted
			continue
		}
		if err = z.Delete("/"+child, true); err != nil {
			return errors.Wrap(err, "Deleting: '/"+child+"'")
		}
	}
	return nil
}

func Test_TriggerPathRefresh(t testing.TB, client Zookeeper) {
	if err := client.refreshPathList(); err != nil {
		t.Fatal("Could not refresh path list in Zookeeper")
	}
}

type Test_InMemoryClient struct {
	root MutableNode
}

func Test_NewInMemoryClient() *Test_InMemoryClient {
	return &Test_InMemoryClient{NewNode("/", []byte{}, []string{})}
}

func (z *Test_InMemoryClient) Close() {
}

func (z *Test_InMemoryClient) SetLogger(_ log.SimpleLogger) {
}

func (z Test_InMemoryClient) Get(path string) ([]byte, Stat, error) {
	node, err := getInHierarchy(z.root, path)
	if err != nil {
		return nil, nil, err
	}
	if nil == node {
		return []byte(""), inMemoryStat{}, nil
	}
	return node.Content(), inMemoryStat{}, nil
}

func getInHierarchy(root Node, path string) (Node, error) {
	if nil == root {
		return nil, errors.New("Root is nil (" + path + ")")
	}
	if root.Path() == path {
		return root, nil
	}

	name, err := childName(root.Path(), path)
	if err != nil {
		return nil, err
	}
	child := root.Child(name)
	if nil == child {
		return nil, nil
	}
	return getInHierarchy(child, path)
}

func toMutableNode(node Node) (MutableNode, error) {
	toReturn, ok := node.(MutableNode)
	if !ok {
		return nil, errors.New("Child node not mutable: " + node.Path())
	}
	return toReturn, nil
}

func (z Test_InMemoryClient) Children(path string) ([]string, Stat, error) {
	node, err := getInHierarchy(z.root, path)
	if err != nil {
		return nil, nil, err
	}
	if nil == node {
		return nil, nil, fmt.Errorf("No node found on path %s", path)
	}
	return node.ChildNames(), inMemoryStat{}, nil
}

func (z *Test_InMemoryClient) Delete(path string, _ int32) error {
	parentNode, err := getInHierarchy(z.root, parent(path))
	if err != nil {
		return err
	}
	elmtName := lastElementName(path)
	child := parentNode.Child(elmtName)
	if nil == child {
		return nil
	}
	if len(child.Childs()) != 0 {
		return zk.ErrNotEmpty
	}

	mNode, err := toMutableNode(parentNode)
	if err != nil {
		return err
	}
	mNode.DeleteChild(elmtName)
	return nil
}

func lastElementName(path string) string {
	splitted := strings.Split(path, "/")
	return splitted[len(splitted)-1]
}

func (z Test_InMemoryClient) Exists(path string) (bool, Stat, error) {
	if nil == z.root && "/" == path {
		return false, inMemoryStat{}, nil
	}
	node, err := getInHierarchy(z.root, path)
	if err != nil {
		return false, nil, err
	}
	return nil != node, inMemoryStat{}, nil
}

func (z *Test_InMemoryClient) Create(path string, data []byte, _ int32, _ []zk.ACL) (string, error) {
	_, err := z.Set(path, data, 0)
	return "", err
}

func (z *Test_InMemoryClient) Set(path string, data []byte, _ int32) (Stat, error) {
	if err := set(z.root, path, data); err != nil {
		return inMemoryStat{}, err
	}
	return inMemoryStat{}, nil
}

func (z *Test_InMemoryClient) initRoot() {
	if nil == z.root {
		z.root = NewDefinedNode("/", []byte{}, []Node{})
	}
}

func set(node Node, path string, content []byte) error {
	if !strings.HasPrefix(path, node.Path()) {
		return fmt.Errorf("Path %s is not parent of %s", node.Path(), path)
	}
	if node.Path() == path {
		mNode, err := toMutableNode(node)
		if err != nil {
			return errors.Wrapf(err, "Shift node to mutable node %s", node.Path())
		}
		mNode.SetContent(content)
		return nil
	}
	name, err := childName(node.Path(), path)
	if err != nil {
		return err
	}
	child := node.Child(name)
	if nil == child {
		child = NewNode(childPath(node.Path(), name), []byte{}, []string{})
		mNode, err := toMutableNode(node)
		if err != nil {
			return errors.Wrapf(err, "Shift node to mutable node %s", node.Path())
		}
		mNode.AddChild(child)
	}
	return set(child, path, content)
}

func childName(parentPath, fullPath string) (string, error) {
	if parentPath == fullPath {
		return "", fmt.Errorf("Parent path %s == Full path", parentPath)
	}
	if isRoot(parentPath) {
		splittedPath := strings.Split(fullPath, "/")
		return splittedPath[1], nil
	}
	withoutParent := strings.Replace(fullPath, parentPath, "", 1)
	splitted := strings.Split(withoutParent, "/")
	return splitted[1], nil
}

func isRoot(path string) bool {
	return "/" == path
}

type inMemoryStat struct {
	Vers int32
}

func (z inMemoryStat) Version() int32 {
	return z.Vers
}
