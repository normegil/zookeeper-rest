package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/normegil/zookeeper-rest/modules/environment"
	errPkg "github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/normegil/zookeeper-rest/modules/test"
	"github.com/normegil/zookeeper-rest/modules/zookeeper"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func TestController(t *testing.T) {
	zkDatas := zookeeper.NewIdentifiedNode("1b011870-db9f-42e3-b03a-3142093d483d", zookeeper.NewDefinedNode("/", []byte("root"), []zookeeper.Node{
		zookeeper.NewIdentifiedNode("7f92685c-0103-4bf2-bd27-7ab96928a037", zookeeper.NewNode("/home", []byte("home"), []string{})),
		zookeeper.NewIdentifiedNode("20fbd7ad-515a-4417-9eda-90ef89cc121f", zookeeper.NewDefinedNode("/var", []byte("var"), []zookeeper.Node{
			zookeeper.NewIdentifiedNode("e1d7e8f6-9778-4f35-b1e0-72b7160cff09", zookeeper.NewNode("/var/test", []byte("{\"test\": \"value\"}"), []string{})),
		})),
		zookeeper.NewIdentifiedNode("205a8193-3c5f-40d0-8021-c53a0a87b37e", zookeeper.NewDefinedNode("/usr", []byte("user123"), []zookeeper.Node{
			zookeeper.NewIdentifiedNode("8b3fb64a-53f4-46f7-b867-354376b872be", zookeeper.NewNode("/usr/bin", []byte("binaries"), []string{})),
			zookeeper.NewIdentifiedNode("4a291277-b5a1-4930-84ff-938d6522452d", zookeeper.NewNode("/usr/share", []byte("shared folder"), []string{})),
			zookeeper.NewIdentifiedNode("814e930e-aa16-4c64-97b6-e1e0deacc1d1", zookeeper.NewNode("/usr/lib", []byte("Libraries"), []string{})),
		},
		)),
	}))

	zkCli, closeZkFn := zookeeper.Test_NewZookeeperClient(t)
	defer closeZkFn()

	c := Controller{environment.Env{Logger: logrus.New(), Zk: zkCli}}
	nodeHandlers := make(map[string]map[string]httprouter.Handle)
	for _, route := range c.Routes() {
		method := route.Method()
		if _, ok := nodeHandlers[method]; !ok {
			nodeHandlers[method] = make(map[string]httprouter.Handle)
		}
		nodeByPath := nodeHandlers[method]
		nodeByPath[route.Path()] = route.Handler()
	}

	t.Run("Get", func(t *testing.T) {
		reset(t, zkCli, zkDatas)
		baseURL := "http://example.com"
		serviceURL := baseURL + "/rest/node"
		configNodes := []string{"/zookeeper", "/zk-rest"}
		configNodeIDs, err := zkCli.IDs(configNodes)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "IDs for config nodes"))
		}
		testcases := []struct {
			Path                string
			ExpectedNodeID      string
			ExpectedNodeContent string
			ExpectedNodePath    string
			ExpectedNodeChilds  map[string]string
			ExpectedStatus      int
			ExpectedErrorCode   int
		}{
			{
				Path:                "/",
				ExpectedStatus:      http.StatusOK,
				ExpectedNodeID:      "1b011870-db9f-42e3-b03a-3142093d483d",
				ExpectedNodePath:    "/",
				ExpectedNodeContent: "root",
				ExpectedNodeChilds: map[string]string{
					"/home":      serviceURL + "/7f92685c-0103-4bf2-bd27-7ab96928a037",
					"/var":       serviceURL + "/20fbd7ad-515a-4417-9eda-90ef89cc121f",
					"/usr":       serviceURL + "/205a8193-3c5f-40d0-8021-c53a0a87b37e",
					"/zk-rest":   serviceURL + "/" + configNodeIDs["/zk-rest"],
					"/zookeeper": serviceURL + "/" + configNodeIDs["/zookeeper"],
				},
			},
			{
				Path:                "",
				ExpectedStatus:      http.StatusOK,
				ExpectedNodeID:      "1b011870-db9f-42e3-b03a-3142093d483d",
				ExpectedNodePath:    "/",
				ExpectedNodeContent: "root",
				ExpectedNodeChilds: map[string]string{
					"/home":      serviceURL + "/7f92685c-0103-4bf2-bd27-7ab96928a037",
					"/var":       serviceURL + "/20fbd7ad-515a-4417-9eda-90ef89cc121f",
					"/usr":       serviceURL + "/205a8193-3c5f-40d0-8021-c53a0a87b37e",
					"/zk-rest":   serviceURL + "/" + configNodeIDs["/zk-rest"],
					"/zookeeper": serviceURL + "/" + configNodeIDs["/zookeeper"],
				},
			},
			{
				Path:                "/home",
				ExpectedStatus:      http.StatusOK,
				ExpectedNodeID:      "7f92685c-0103-4bf2-bd27-7ab96928a037",
				ExpectedNodePath:    "/home",
				ExpectedNodeContent: "home",
				ExpectedNodeChilds:  map[string]string{},
			},
			{
				Path:                "/usr/bin",
				ExpectedStatus:      http.StatusOK,
				ExpectedNodeID:      "8b3fb64a-53f4-46f7-b867-354376b872be",
				ExpectedNodePath:    "/usr/bin",
				ExpectedNodeContent: "binaries",
				ExpectedNodeChilds:  map[string]string{},
			},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				id, err := zkCli.ID(testdata.Path)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Getting %s ID", testdata.Path))
				}
				request := httptest.NewRequest("GET", fmt.Sprintf(serviceURL+"/%s", id), nil)
				result := httptest.NewRecorder()
				handlers(t, nodeHandlers, GET_METHOD, GET_PATH)(result, request, httprouter.Params([]httprouter.Param{{
					Key:   NODE_ID_PARAM_KEY,
					Value: id,
				}}))
				if testdata.ExpectedStatus != result.Code {
					t.Error(test.Format("Expected response status is not equals to received status", testdata.ExpectedStatus, result.Code))
				}
				if testdata.ExpectedStatus != http.StatusOK {
					var response errPkg.ErrorResponse
					if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
						t.Fatal(errors.Wrapf(err, "Could not unmarshal response body %s", string(result.Body.Bytes())))
					}
					if testdata.ExpectedErrorCode != response.Code {
						t.Fatal(test.Format("Expected error code  doesnt match response code", testdata.ExpectedErrorCode, strconv.Itoa(response.Code)+": "+response.Err.Error()))
					}
				} else {
					var response nodeResponse
					if err := json.Unmarshal(result.Body.Bytes(), &response); nil != err {
						t.Fatal(errors.Wrapf(err, "Could not unmarshal response"))
					}
					if testdata.ExpectedNodeID != response.ID {
						t.Error(test.Format("Expected ID is not the node ID", testdata.ExpectedNodeID, response.ID))
					}
					if testdata.ExpectedNodePath != response.Path {
						t.Error(test.Format("Expected path is not the node path", testdata.Path, response.Path))
					}
					if testdata.ExpectedNodeContent != response.Content {
						t.Error(test.Format("Expected content is not the node content", testdata.ExpectedNodeContent, response.Content))
					}
					expectedURL := fmt.Sprintf(serviceURL+"/%s", id)
					if expectedURL != response.URL.String() {
						t.Error(test.Format("Expected content is not the node content", testdata.ExpectedNodeContent, response.Content))
					}
					if len(testdata.ExpectedNodeChilds) != len(response.Childs) {
						t.Error(test.Format("Expected childs length is not the returned child length", testdata.ExpectedNodeChilds, response.Childs))
					}
					for path, childURLStr := range testdata.ExpectedNodeChilds {
						childToTest, found := response.Childs[path]
						if !found {
							t.Errorf("Child URL not found (%s). Expected: %s", path, childURLStr)
						} else if childToTest.String() != childURLStr {
							t.Error(test.Format("Child URL not equals expected URL ("+path+")", childURLStr, childToTest))
						}
					}
				}
			})
		}
	})

	t.Run("Create", func(t *testing.T) {
		testcases := []struct {
			Path              string
			Content           string
			ExpectedStatus    int
			ExpectedErrorCode int
		}{
			{"/opt", "opt", http.StatusOK, 0},
			{"/opt/whatever", "whatever", http.StatusOK, 0},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				reset(t, zkCli, zkDatas)

				body := map[string]string{
					"path":    testdata.Path,
					"content": testdata.Content,
				}
				marshaledBody, err := json.Marshal(body)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Marshalling body %+v", body))
				}
				r := httptest.NewRequest(CREATE_METHOD, CREATE_PATH, bytes.NewReader(marshaledBody))
				result := httptest.NewRecorder()
				handlers(t, nodeHandlers, CREATE_METHOD, CREATE_PATH)(result, r, httprouter.Params([]httprouter.Param{}))

				if testdata.ExpectedStatus != result.Code {
					t.Error(test.Format("Expected response status is not equals to received status", testdata.ExpectedStatus, result.Code))
				}
				if testdata.ExpectedStatus != http.StatusOK {
					var response errPkg.ErrorResponse
					if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
						t.Fatal(errors.Wrapf(err, "Could not unmarshal response body %s", string(result.Body.Bytes())))
					}
					if testdata.ExpectedErrorCode != response.Code {
						t.Fatal(test.Format("Expected error code  doesnt match response code", testdata.ExpectedErrorCode, strconv.Itoa(response.Code)+": "+response.Err.Error()))
					}
				} else {
					exist, err := zkCli.Exist(testdata.Path)
					if err != nil {
						t.Fatal(errors.Wrapf(err, "Check %s existence", testdata.Path))
					}
					if !exist {
						t.Error(fmt.Sprintf("Node doesn't exist ('%s')", testdata.Path))
					}
					content, err := zkCli.Content(testdata.Path)
					if nil != err {
						t.Fatal(errors.Wrapf(err, "Getting content of %s", testdata.Path))
					}
					if testdata.Content != string(content) {
						t.Error(test.Format("Setted Content not equal to sent content", testdata.Content, string(content)))
					}
				}
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		testcases := []struct {
			Path              string
			ExpectedStatus    int
			ExpectedErrorCode int
		}{
			{"/usr", http.StatusOK, 0},
			{"/usr/bin", http.StatusOK, 0},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"'", func(t *testing.T) {
				reset(t, zkCli, zkDatas)

				content := uuid.NewV4().String()
				body := map[string]string{
					"content": content,
				}
				marshaledBody, err := json.Marshal(body)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Marshalling body %+v", body))
				}
				id, err := zkCli.ID(testdata.Path)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Getting ID of %s", testdata.Path))
				}
				baseURL := "http://example.com"
				serviceURL := baseURL + "/rest/node"
				fullURL := fmt.Sprintf(serviceURL+"/%s", id)
				r := httptest.NewRequest(UPDATE_METHOD, fullURL, bytes.NewReader(marshaledBody))
				result := httptest.NewRecorder()
				handlers(t, nodeHandlers, UPDATE_METHOD, UPDATE_PATH)(result, r, httprouter.Params([]httprouter.Param{{
					Key:   NODE_ID_PARAM_KEY,
					Value: id,
				}}))

				if testdata.ExpectedStatus != result.Code {
					t.Error(test.Format("Expected response status is not equals to received status", testdata.ExpectedStatus, result.Code))
				}
				if testdata.ExpectedStatus != http.StatusOK {
					var response errPkg.ErrorResponse
					if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
						t.Fatal(errors.Wrapf(err, "Could not unmarshal response body %s", string(result.Body.Bytes())))
					}
					if testdata.ExpectedErrorCode != response.Code {
						t.Fatal(test.Format("Expected error code  doesnt match response code", testdata.ExpectedErrorCode, strconv.Itoa(response.Code)+": "+response.Err.Error()))
					}
				} else {
					exist, err := zkCli.Exist(testdata.Path)
					if err != nil {
						t.Fatal(errors.Wrapf(err, "Check %s existence", testdata.Path))
					}
					if !exist {
						t.Error(fmt.Sprintf("Node doesn't exist ('%s')", testdata.Path))
					}
					settedContent, err := zkCli.Content(testdata.Path)
					if nil != err {
						t.Fatal(errors.Wrapf(err, "Getting content of %s", testdata.Path))
					}
					if content != string(settedContent) {
						t.Error(test.Format("Setted Content not equal to sent content", string(settedContent), content))
					}
				}
			})
		}
	})

	t.Run("Delete", func(t *testing.T) {
		testcases := []struct {
			Path              string
			Recursive         string
			ExpectedStatus    int
			ExpectedErrorCode int
		}{
			{"/home", "", http.StatusOK, 0},
			{"/usr/bin", "", http.StatusOK, 0},
			{"/usr", "true", http.StatusOK, 0},
			{"/usr", "false", http.StatusBadRequest, 40005},
		}
		for _, testdata := range testcases {
			t.Run("Path:'"+testdata.Path+"';Recurse:"+testdata.Recursive, func(t *testing.T) {
				reset(t, zkCli, zkDatas)
				id, err := zkCli.ID(testdata.Path)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "Getting %s ID", testdata.Path))
				}
				url := fmt.Sprintf("http://example.com/rest/node/%s?recursive=%s", id, testdata.Recursive)
				r := httptest.NewRequest("DELETE", url, nil)
				result := httptest.NewRecorder()
				handlers(t, nodeHandlers, DELETE_METHOD, DELETE_PATH)(result, r, httprouter.Params([]httprouter.Param{{
					Key:   NODE_ID_PARAM_KEY,
					Value: id,
				}}))
				if testdata.ExpectedStatus != result.Code {
					t.Error(test.Format("Expected response status is not equals to received status", testdata.ExpectedStatus, result.Code))
				}
				if testdata.ExpectedStatus != http.StatusOK {
					var response errPkg.ErrorResponse
					if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
						t.Fatal(errors.Wrapf(err, "Could not unmarshal response body %s", string(result.Body.Bytes())))
					}
					if testdata.ExpectedErrorCode != response.Code {
						t.Fatal(test.Format("Expected error code  doesnt match response code", testdata.ExpectedErrorCode, strconv.Itoa(response.Code)+": "+response.Err.Error()))
					}
				} else {
					exist, err := zkCli.Exist(testdata.Path)
					if err != nil {
						t.Fatal(errors.Wrapf(err, "Check %s existence", testdata.Path))
					}
					if exist {
						t.Error(fmt.Sprintf("Node still exist ('%s')", testdata.Path))
					}
				}
			})
		}
	})
}

func reset(t testing.TB, cli zookeeper.Zookeeper, zkData zookeeper.Node) {
	zookeeper.Test_CleanInsert(t, cli, zkData)
}

func handlers(t testing.TB, allHandlers map[string]map[string]httprouter.Handle, method, path string) httprouter.Handle {
	handlerByPath, ok := allHandlers[method]
	if !ok {
		t.Fatal("Could not find method %s in handlers", method)
	}
	handler, ok := handlerByPath[path]
	if !ok {
		t.Fatalf("Could not find handler for method %s and path %s. Found: %+v", method, path, handlerByPath)
	}
	return handler
}

func contains(slice []string, toTest string) bool {
	for _, str := range slice {
		if str == toTest {
			return true
		}
	}
	return false
}
