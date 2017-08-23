package zookeeper

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Identified interface {
	ID() string
}

type Node interface {
	Path() string
	Content() []byte
	IsRecursive() bool
	Equals(Node) bool
	EqualsRecursive(Node) bool
	ChildNames() []string
	Childs() []Node
	Child(string) Node
}

type MutableNode interface {
	Node
	SetContent([]byte)
	AddChild(Node)
	DeleteChild(string)
}

type nodeImpl struct {
	path      string
	content   []byte
	childs    []Node
	recursive bool
}

type prefixedStringer interface {
	PrefixedString(prefix string) string
}

func NewNode(path string, content []byte, childs []string) *nodeImpl {
	childNodes := make([]Node, len(childs))
	basePath := path
	if path != "/" {
		basePath += "/"
	}
	for i, childName := range childs {
		childNodes[i] = &nodeImpl{path: basePath + childName}
	}
	return NewDefinedNode(path, content, childNodes)
}

func NewRecursiveNode(path string, z Zookeeper) (*nodeImpl, error) {
	childs, err := z.Children(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Loading children of %s", path)
	}
	childNodes := make([]Node, len(childs))
	for i, childName := range childs {
		childPath := childPath(path, childName)
		childNodes[i], err = NewRecursiveNode(childPath, z)
		if err != nil {
			return nil, err
		}
	}

	content, err := z.Content(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Loading content of %s", path)
	}
	node := NewDefinedNode(path, content, childNodes)
	node.recursive = true
	return node, nil
}

func NewDefinedNode(path string, content []byte, childs []Node) *nodeImpl {
	return &nodeImpl{
		path:    path,
		content: content,
		childs:  childs,
	}
}

func (n nodeImpl) Path() string {
	return n.path
}

func (n nodeImpl) Content() []byte {
	return n.content
}

func (n nodeImpl) IsRecursive() bool {
	return n.recursive
}

func (n nodeImpl) ChildNames() []string {
	childs := n.Childs()
	names := make([]string, len(childs))
	for i, child := range childs {
		names[i] = lastPathElementName(child.Path())
	}
	return names
}

func (n nodeImpl) Childs() []Node {
	toReturn := make([]Node, len(n.childs))
	for i, child := range n.childs {
		toReturn[i] = child
	}
	return toReturn
}

func (n nodeImpl) Equals(other Node) bool {
	if n.Path() != other.Path() {
		return false
	}

	if len(n.Content()) != len(other.Content()) {
		return false
	}

	if len(n.Childs()) != len(other.Childs()) {
		return false
	}

	for i, nByte := range n.Content() {
		if nByte != other.Content()[i] {
			return false
		}
	}

	for _, nChildName := range n.ChildNames() {
		found := false
		for _, oChildName := range other.ChildNames() {
			if oChildName == nChildName {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (n nodeImpl) EqualsRecursive(other Node) bool {
	if !n.Equals(other) {
		return false
	}

	for _, nChild := range n.Childs() {
		for _, oChild := range other.Childs() {
			if nChild.Path() == oChild.Path() {
				if nChild.EqualsRecursive(oChild) {
					return false
				}
			}
		}
	}
	return true
}

func (n *nodeImpl) SetContent(data []byte) {
	n.content = data
}

func (n nodeImpl) String() string {
	return n.PrefixedString("")
}

func (n nodeImpl) PrefixedString(prefix string) string {
	base := fmt.Sprintf(prefix+"{Path:%s;Content:%s;Childs:", n.Path(), n.Content())
	if !n.IsRecursive() {
		return base + fmt.Sprintf("%+v}", n.ChildNames())
	}
	for _, child := range n.Childs() {
		if childPrefixedStringer, ok := child.(prefixedStringer); ok {
			base += "\n" + childPrefixedStringer.PrefixedString(prefix+"\t")
		} else if childStringer, ok := child.(fmt.Stringer); ok {
			base += "\n" + childStringer.String()
		} else {
			base += fmt.Sprintf("\n%+v", child)
		}
	}
	return base + "}"
}

func (n nodeImpl) Child(childName string) Node {
	for _, child := range n.Childs() {
		if strings.HasSuffix(child.Path(), "/"+childName) {
			return child
		}
	}
	return nil
}

func (n *nodeImpl) AddChild(child Node) {
	n.childs = append(n.childs, child)
}

func (n *nodeImpl) DeleteChild(childName string) {
	for i, child := range n.childs {
		splittedPath := strings.Split(child.Path(), "/")
		if splittedPath[len(splittedPath)-1] == childName {
			n.childs = append(n.childs[:i], n.childs[i+1:]...)
			break
		}
	}
}

func childPath(parent, childName string) string {
	if "/" == parent {
		return "/" + childName
	}
	return parent + "/" + childName
}

func lastPathElementName(path string) string {
	splittedPath := strings.Split(path, "/")
	return splittedPath[len(splittedPath)-1]
}

type IdentifiedNodeImpl struct {
	id string
	Node
}

func NewIdentifiedNode(id string, associated Node) *IdentifiedNodeImpl {
	return &IdentifiedNodeImpl{
		id:   id,
		Node: associated,
	}
}

func (i IdentifiedNodeImpl) ID() string {
	return i.id
}
