package zookeeper_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/normegil/zookeeper-rest/modules/test"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func TestWithZookeeper(t *testing.T) {
	client, closeClientFn := zookeeper.Test_NewZookeeperClient(t)
	defer closeClientFn()
	testRoot := zookeeper.NewDefinedNode("/", []byte("root"), []zookeeper.Node{
		zookeeper.NewDefinedNode("/home", []byte("home"), []zookeeper.Node{}),
		zookeeper.NewDefinedNode("/var", []byte("var"), []zookeeper.Node{
			zookeeper.NewDefinedNode("/var/test", []byte("{\"test\": \"value\"}"), []zookeeper.Node{}),
		}),
		zookeeper.NewDefinedNode("/usr", []byte("user123"), []zookeeper.Node{
			zookeeper.NewDefinedNode("/usr/bin", []byte("binaries"), []zookeeper.Node{}),
			zookeeper.NewDefinedNode("/usr/share", []byte("shared folder"), []zookeeper.Node{}),
			zookeeper.NewDefinedNode("/usr/lib", []byte("Libraries"), []zookeeper.Node{}),
		}),
	})

	t.Run("Exist", func(t *testing.T) {
		zookeeper.Test_CleanInsert(t, client, testRoot)
		testcases := []struct {
			Path     string
			Expected bool
		}{
			{"/", true},
			{"/opt", false},
			{"/opt/discord", false},
			{"/home", true},
			{"/var", true},
			{"/var/cache", false},
			{"/var/test", true},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				exist, err := client.Exist(testdata.Path)
				if err != nil {
					t.Fatal(errors.Wrap(err, "Getting testing content"))
				}

				if testdata.Expected != exist {
					t.Error(test.Format("Existence check failed", testdata.Expected, exist))
				}
			})
		}
	})

	t.Run("Content", func(t *testing.T) {
		zookeeper.Test_CleanInsert(t, client, testRoot)
		testcases := []struct {
			Path     string
			Expected string
		}{
			{"/", "root"},
			{"/home", "home"},
			{"/var", "var"},
			{"/var/test", "{\"test\": \"value\"}"},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				content, err := client.Content(testdata.Path)
				if err != nil {
					t.Fatal(errors.Wrap(err, "Getting testing content"))
				}

				if testdata.Expected != string(content) {
					t.Error(test.Format("Contents differ", testdata.Expected, content))
				}
			})
		}
	})

	t.Run("Children", func(t *testing.T) {
		zookeeper.Test_CleanInsert(t, client, testRoot)
		testcases := []struct {
			Path           string
			ExpectedChilds []string
		}{
			{Path: "/", ExpectedChilds: []string{"home", "var", "usr", "zookeeper"}},
			{Path: "/home", ExpectedChilds: []string{}},
			{Path: "/usr", ExpectedChilds: []string{"bin", "share", "lib"}},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				childs, err := client.Children(testdata.Path)
				if err != nil {
					t.Fatal(errors.Wrap(err, "Getting children of "+testdata.Path))
				}
				if len(testdata.ExpectedChilds) != len(childs) {
					t.Fatal(test.Format("Length: Childs differs from expected results", testdata.ExpectedChilds, childs))
				}
				for _, expected := range testdata.ExpectedChilds {
					found := false
					for _, child := range childs {
						if child == expected {
							found = true
						}
					}
					if !found {
						t.Errorf("Expected Child not found (%s) in: \n\t%+v", expected, childs)
					}
				}
			})
		}
	})

	t.Run("Load", func(t *testing.T) {
		zookeeper.Test_CleanInsert(t, client, testRoot)
		testcases := []struct {
			Path     string
			Expected zookeeper.Node
		}{
			{"/", *zookeeper.NewNode("/", []byte("root"), []string{"zookeeper", "home", "var", "usr"})},
			{"/home", *zookeeper.NewNode("/home", []byte("home"), []string{})},
			{"/var/test", *zookeeper.NewNode("/var/test", []byte("{\"test\": \"value\"}"), []string{})},
			{"/usr", *zookeeper.NewNode("/usr", []byte("user123"), []string{"bin", "share", "lib"})},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				node, err := client.Load(testdata.Path, false)
				if err != nil {
					t.Fatal(errors.Wrap(err, "Getting testing content"))
				}

				if !testdata.Expected.Equals(node) {
					t.Error(test.Format("Contents differ", testdata.Expected, node))
				}
			})
		}
	})

	t.Run("Load-Recursive", func(t *testing.T) {
		zookeeper.Test_CleanInsert(t, client, testRoot)
		zookeeperNode, err := client.Load("/zookeeper", true)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "Loading /zookeeper"))
		}
		testcases := []struct {
			Path     string
			Expected zookeeper.Node
		}{
			{"/", zookeeper.NewDefinedNode("/", []byte("root"), []zookeeper.Node{
				zookeeper.NewDefinedNode("/home", []byte("home"), []zookeeper.Node{}),
				zookeeper.NewDefinedNode("/var", []byte("var"), []zookeeper.Node{
					zookeeper.NewDefinedNode("/var/test", []byte("{\"test\": \"value\"}"), []zookeeper.Node{}),
				}),
				zookeeper.NewDefinedNode("/usr", []byte("user123"), []zookeeper.Node{
					zookeeper.NewDefinedNode("/usr/bin", []byte("binaries"), []zookeeper.Node{}),
					zookeeper.NewDefinedNode("/usr/share", []byte("shared folder"), []zookeeper.Node{}),
					zookeeper.NewDefinedNode("/usr/lib", []byte("Libraries"), []zookeeper.Node{}),
				}),
				zookeeperNode,
			})},
			{"/home", *zookeeper.NewNode("/home", []byte("home"), []string{})},
			{"/var/test", *zookeeper.NewNode("/var/test", []byte("{\"test\": \"value\"}"), []string{})},
			{"/usr", *zookeeper.NewDefinedNode("/usr", []byte("user123"), []zookeeper.Node{
				zookeeper.NewDefinedNode("bin", []byte(""), []zookeeper.Node{}),
				zookeeper.NewDefinedNode("share", []byte(""), []zookeeper.Node{}),
				zookeeper.NewDefinedNode("lib", []byte(""), []zookeeper.Node{}),
			})},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				node, err := client.Load(testdata.Path, false)
				if err != nil {
					t.Fatal(errors.Wrap(err, "Getting testing content"))
				}

				if !testdata.Expected.Equals(node) {
					t.Error(test.Format("Contents differ", testdata.Expected, node))
				}
			})
		}
	})

	t.Run("Set", func(t *testing.T) {
		testcases := []struct {
			Path     string
			Content  []byte
			Expected zookeeper.Node
		}{
			{
				Path:     "/",
				Content:  []byte("ROOT"),
				Expected: *zookeeper.NewNode("/", []byte("ROOT"), []string{"zookeeper", "home", "var", "usr"}),
			},
			{
				Path:     "/opt",
				Content:  []byte("opt"),
				Expected: *zookeeper.NewNode("/opt", []byte("opt"), []string{}),
			},
			{
				Path:     "/usr/bin",
				Content:  []byte("bin123"),
				Expected: *zookeeper.NewNode("/usr/bin", []byte("bin123"), []string{}),
			},
			{
				Path:     "/opt/whatever",
				Content:  []byte("whatever123"),
				Expected: *zookeeper.NewNode("/opt/whatever", []byte("whatever123"), []string{}),
			},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				zookeeper.Test_CleanInsert(t, client, testRoot)
				if err := client.Set(testdata.Path, testdata.Content); err != nil {
					t.Fatal(errors.Wrap(err, "Setting node content"))
				}

				node, err := client.Load(testdata.Path, false)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Getting tested node content"))
				}

				if !testdata.Expected.Equals(node) {
					t.Error(test.Format("Contents differ", testdata.Expected, node))
				}
			})
		}
	})

	t.Run("Delete", func(t *testing.T) {
		zookeeperNode, err := client.Load("/zookeeper", true)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "Loading /zookeeper"))
		}
		testcases := []struct {
			Root              zookeeper.Node
			Path              string
			ExpectedStructure zookeeper.Node
		}{
			{
				Path: "/home",
				Root: zookeeper.NewDefinedNode("/", []byte{}, []zookeeper.Node{
					zookeeper.NewDefinedNode("/home", []byte{}, []zookeeper.Node{}),
					zookeeper.NewDefinedNode("/var", []byte{}, []zookeeper.Node{}),
				}),
				ExpectedStructure: zookeeper.NewDefinedNode("/", []byte{}, []zookeeper.Node{
					zookeeper.NewDefinedNode("/var", []byte{}, []zookeeper.Node{}),
					zookeeperNode,
				}),
			},
			{
				Path: "/usr/bin",
				Root: zookeeper.NewDefinedNode("/", []byte{}, []zookeeper.Node{
					zookeeper.NewDefinedNode("/var", []byte{}, []zookeeper.Node{}),
					zookeeper.NewDefinedNode("/usr", []byte{}, []zookeeper.Node{
						zookeeper.NewDefinedNode("/usr/bin", []byte{}, []zookeeper.Node{}),
						zookeeper.NewDefinedNode("/usr/share", []byte{}, []zookeeper.Node{}),
					}),
				}),
				ExpectedStructure: zookeeper.NewDefinedNode("/", []byte{}, []zookeeper.Node{
					zookeeper.NewDefinedNode("/var", []byte{}, []zookeeper.Node{}),
					zookeeper.NewDefinedNode("/usr", []byte{}, []zookeeper.Node{
						zookeeper.NewDefinedNode("/usr/share", []byte{}, []zookeeper.Node{}),
					}),
					zookeeperNode,
				}),
			},
			{
				Path: "/usr",
				Root: zookeeper.NewDefinedNode("/", []byte{}, []zookeeper.Node{
					zookeeper.NewDefinedNode("/var", []byte{}, []zookeeper.Node{}),
					zookeeper.NewDefinedNode("/usr", []byte{}, []zookeeper.Node{
						zookeeper.NewDefinedNode("/usr/bin", []byte{}, []zookeeper.Node{}),
						zookeeper.NewDefinedNode("/usr/share", []byte{}, []zookeeper.Node{}),
					}),
				}),
				ExpectedStructure: zookeeper.NewDefinedNode("/", []byte{}, []zookeeper.Node{
					zookeeper.NewDefinedNode("/var", []byte{}, []zookeeper.Node{}),
					zookeeperNode,
				}),
			},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				zookeeper.Test_CleanInsert(t, client, testdata.Root)
				if err := client.Delete(testdata.Path, true); err != nil {
					t.Fatal(errors.Wrap(err, "Deleting: '"+testdata.Path+"'"))
				}

				expectedCli := zookeeper.Zookeeper{
					Logger:             logrus.New(),
					SubstitutionClient: zookeeper.Test_NewInMemoryClient(),
				}
				if err := zookeeper.Test_InsertNode(expectedCli, testdata.ExpectedStructure); err != nil {
					t.Fatal(errors.Wrapf(err, "Inserting expected data into associated client"))
				}
				currentNode, err := client.Load("/", true)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Loading current '/'"))
				}
				if !testdata.ExpectedStructure.EqualsRecursive(currentNode) {
					test.Format("Structure not equal", testdata.ExpectedStructure, currentNode)
				}
			})
		}
	})

	t.Run("Path", func(t *testing.T) {
		zookeeper.Test_CleanInsert(t, client, testRoot)
		pathAssociations := make(map[string]string)
		pathAssociations["f051edfe-4dc7-4d37-b732-feb01cfc7dd1"] = "/"
		pathAssociations["22748fcb-d876-48a5-9836-fd01e61d03b3"] = "/home"
		pathAssociations["ab360a73-d7e1-4a17-a0be-fe9640961ea8"] = "/usr"
		pathAssociations["8c494086-b018-4be5-ba68-033ddf0a9411"] = "/usr/bin"
		content, err := json.Marshal(pathAssociations)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "Marhal content of Zookeeper paths configurations node"))
		}
		if err := client.Set(zookeeper.ZOOKEEPER_NODE_PATH_STORAGE, content); err != nil {
			t.Fatal(errors.Wrapf(err, "Set content of %s", zookeeper.ZOOKEEPER_NODE_PATH_STORAGE))
		}
		testcases := []struct {
			ID       string
			Expected string
		}{
			{"f051edfe-4dc7-4d37-b732-feb01cfc7dd1", "/"},
			{"22748fcb-d876-48a5-9836-fd01e61d03b3", "/home"},
			{"ab360a73-d7e1-4a17-a0be-fe9640961ea8", "/usr"},
			{"8c494086-b018-4be5-ba68-033ddf0a9411", "/usr/bin"},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Expected+"'", func(t *testing.T) {
				path, err := client.Path(testdata.ID)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Loading path for "+testdata.ID))
				}
				if testdata.Expected != path {
					t.Error(test.Format("Path returned is not the expected path", testdata.Expected, path))
				}
			})
		}
	})

	t.Run("ID", func(t *testing.T) {
		zookeeper.Test_CleanInsert(t, client, testRoot)
		pathAssociations := make(map[string]string)
		pathAssociations["f051edfe-4dc7-4d37-b732-feb01cfc7dd1"] = "/"
		pathAssociations["22748fcb-d876-48a5-9836-fd01e61d03b3"] = "/home"
		pathAssociations["ab360a73-d7e1-4a17-a0be-fe9640961ea8"] = "/usr"
		pathAssociations["8c494086-b018-4be5-ba68-033ddf0a9411"] = "/usr/bin"
		content, err := json.Marshal(pathAssociations)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "Marhal content of Zookeeper paths configurations node"))
		}
		if err := client.Set(zookeeper.ZOOKEEPER_NODE_PATH_STORAGE, content); err != nil {
			t.Fatal(errors.Wrapf(err, "Set content of %s", zookeeper.ZOOKEEPER_NODE_PATH_STORAGE))
		}
		testcases := []struct {
			Path     string
			Expected string
		}{
			{"/", "f051edfe-4dc7-4d37-b732-feb01cfc7dd1"},
			{"/home", "22748fcb-d876-48a5-9836-fd01e61d03b3"},
			{"/usr", "ab360a73-d7e1-4a17-a0be-fe9640961ea8"},
			{"/usr/bin", "8c494086-b018-4be5-ba68-033ddf0a9411"},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				id, err := client.ID(testdata.Path)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Loading id for "+testdata.Path))
				}
				if testdata.Expected != id {
					t.Error(test.Format("ID returned is not the expected ID", testdata.Expected, id))
				}
			})
		}
	})

	t.Run("IDs", func(t *testing.T) {
		zookeeper.Test_CleanInsert(t, client, testRoot)
		pathAssociations := make(map[string]string)
		pathAssociations["f051edfe-4dc7-4d37-b732-feb01cfc7dd1"] = "/"
		pathAssociations["22748fcb-d876-48a5-9836-fd01e61d03b3"] = "/home"
		pathAssociations["ab360a73-d7e1-4a17-a0be-fe9640961ea8"] = "/usr"
		pathAssociations["8c494086-b018-4be5-ba68-033ddf0a9411"] = "/usr/bin"
		content, err := json.Marshal(pathAssociations)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "Marhal content of Zookeeper paths configurations node"))
		}
		if err := client.Set(zookeeper.ZOOKEEPER_NODE_PATH_STORAGE, content); err != nil {
			t.Fatal(errors.Wrapf(err, "Set content of %s", zookeeper.ZOOKEEPER_NODE_PATH_STORAGE))
		}
		testcases := []struct {
			Paths    []string
			Expected map[string]string
		}{
			{
				Paths: []string{"/", "/home", "/usr", "/usr/bin"},
				Expected: map[string]string{
					"/":        "f051edfe-4dc7-4d37-b732-feb01cfc7dd1",
					"/home":    "22748fcb-d876-48a5-9836-fd01e61d03b3",
					"/usr":     "ab360a73-d7e1-4a17-a0be-fe9640961ea8",
					"/usr/bin": "8c494086-b018-4be5-ba68-033ddf0a9411",
				},
			},
		}
		for _, testdata := range testcases {
			t.Run("", func(t *testing.T) {
				idAssociations, err := client.IDs(testdata.Paths)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Loading ids for %+v", testdata.Paths))
				}
				if len(testdata.Expected) != len(idAssociations) {
					t.Error(test.Format("IDs returned are not the expected IDs (Length problem)", testdata.Expected, idAssociations))
				}
				for path, expectedID := range testdata.Expected {
					if expectedID != idAssociations[path] {
						t.Error(test.Format("IDs returned are not the expected IDs", testdata.Expected, idAssociations))
					}
				}
			})
		}
	})
}

func checkContent(t testing.TB, path string, expected zookeeper.Zookeeper, toTest zookeeper.Zookeeper) {
	expectedContent, err := expected.Content(path)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "Getting content of Client 1 on %s", path))
	}
	toTestContent, err := toTest.Content(path)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "Getting content of Client 2 on %s", path))
	}
	if !bytes.Equal(expectedContent, toTestContent) {
		t.Fatal(test.Format("Contents doesn't correspond ("+path+")", expectedContent, toTestContent))
	}
}
