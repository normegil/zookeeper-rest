package zookeeper_test

import (
	"testing"

	"github.com/normegil/zookeeper-rest/modules/test"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
)

func TestChilds(t *testing.T) {
	testcases := []struct {
		Input  zookeeper.Node
		Output []zookeeper.Node
	}{
		{
			Input: zookeeper.NewNode("/", []byte("root"), []string{"usr", "home"}),
			Output: []zookeeper.Node{
				zookeeper.NewDefinedNode("/usr", []byte{}, []zookeeper.Node{}),
				zookeeper.NewDefinedNode("/home", []byte{}, []zookeeper.Node{}),
			},
		},
		{
			Input: zookeeper.NewDefinedNode("/", []byte("root"), []zookeeper.Node{
				zookeeper.NewDefinedNode("/defined1", []byte{}, []zookeeper.Node{}),
				zookeeper.NewDefinedNode("/defined2", []byte{}, []zookeeper.Node{}),
			}),
			Output: []zookeeper.Node{
				zookeeper.NewDefinedNode("/defined1", []byte{}, []zookeeper.Node{}),
				zookeeper.NewDefinedNode("/defined2", []byte{}, []zookeeper.Node{}),
			},
		},
	}
	for _, testdata := range testcases {
		t.Run("Path:'"+testdata.Input.Path()+"'", func(t *testing.T) {
			childs := testdata.Input.Childs()
			t.Logf("Childs: %+v", childs)
			if len(testdata.Output) != len(childs) {
				t.Errorf(test.Format("Lengh of expected childs doesn't correspond to lenght of returned childs", testdata.Output, childs))
			}
			for _, expected := range testdata.Output {
				found := false
				var foundChild zookeeper.Node
				for _, child := range childs {
					if child.Path() == expected.Path() {
						found = true
						foundChild = child
						break
					}
				}
				if !found {
					t.Errorf("Node not found: %s", expected.Path())
				} else if !expected.EqualsRecursive(foundChild) {
					t.Errorf(test.Format("Childs ("+expected.Path()+") not equals", expected, foundChild))
				}
			}
		})
	}
}
