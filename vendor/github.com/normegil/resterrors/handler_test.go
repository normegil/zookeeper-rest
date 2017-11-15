package resterrors

import (
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/pkg/errors"
)

func TestHandle(t *testing.T) {
	defaultCode := 50000
	definitions := []ErrorDefinition{
		{50000, 500, "http://example.com/wiki/errors/50000", "An undetermined error happened on the Server."},
		{50301, 503, "http://example.com/wiki/errors/50301", "Could not connect to Zookeeper at given address. Check that zookeeper is running and accessible by the rest server."},
		{40001, 400, "http://example.com/wiki/errors/40001", "The request doesn't correspond to the structure needed to solve your request. Please review the body of your request."},
		{40002, 400, "http://example.com/wiki/errors/40002", "A value is missing from your request or is misplaced."},
		{40003, 400, "http://example.com/wiki/errors/40003", "One of your parameter doesn't have the expected format and cannot be parsed."},
		{40004, 400, "http://example.com/wiki/errors/40004", "Cannot remove root path from zookeeper."},
		{40005, 400, "http://example.com/wiki/errors/40005", "Cannot remove a node with existing childs (Use 'recursive' option)."},
		{40100, 401, "http://example.com/wiki/errors/40100", "All authentication headers were detected as empty."},
		{40101, 401, "http://example.com/wiki/errors/40101", "Error while trying to authenticate user using Authorization header."},
		{40102, 401, "http://example.com/wiki/errors/40102", "No error happened but user is empty."},
		{40103, 401, "http://example.com/wiki/errors/40103", "Authentication header didn't contain any user."},
	}
	tests := []struct {
		err     error
		errCode int
	}{
		{errors.New("TestError"), defaultCode},
		{NewErrWithCode(50301, errors.New("TestError")), 50301},
		{NewErrWithCode(40001, errors.New("TestError")), 40001},
		{NewErrWithCode(40002, errors.New("TestError")), 40002},
		{NewErrWithCode(40100, errors.New("TestError")), 40100},
		{NewErrWithCode(40101, errors.New("TestError")), 40101},
		{NewErrWithCode(40102, errors.New("TestError")), 40102},
	}
	for index, testdata := range tests {
		t.Run(testdata.err.Error(), func(t *testing.T) {
			resp := httptest.NewRecorder()

			Handler{definitions, defaultCode}.Handle(resp, testdata.err)

			var expectedErr *ErrorResponse
			for _, definition := range definitions {
				if testdata.errCode == definition.Code {
					var err error
					expectedErr, err = definition.ToResponse(testdata.err)
					if err != nil {
						t.Fatal(errors.Wrapf(err, "Get Response from definition"))
					}
					break
				}
			}
			if expectedErr.HTTPStatus != resp.Code {
				t.Errorf("%d:HTTP Status (%+v) doesn't correspond to expected (%+v).", index, resp.Code, expectedErr.HTTPStatus)
			}

			var body ErrorResponse
			json.Unmarshal(resp.Body.Bytes(), &body)
			expectedErr.Time = body.Time
			if expectedErr.Code != body.Code {
				t.Errorf("Body (%s) doesn't correspond to expected (%s)", strconv.Itoa(body.Code), strconv.Itoa(testdata.errCode))
			}
		})
	}
}
