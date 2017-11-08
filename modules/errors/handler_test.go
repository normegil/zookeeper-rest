package errors

import (
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	errorFormat "github.com/normegil/formats/error"
	"github.com/normegil/zookeeper-rest/modules/test"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func TestHandle(t *testing.T) {
	errResponses := loadErrorRessources("assets/errors.csv")
	tests := []struct {
		err     error
		errCode int
	}{
		{errors.New("TestError"), 50000},
		{NewErrWithCode(50301, errors.New("TestError")), 50301},
		{NewErrWithCode(40001, errors.New("TestError")), 40001},
		{NewErrWithCode(40002, errors.New("TestError")), 40002},
		{NewErrWithCode(40100, errors.New("TestError")), 40100},
		{NewErrWithCode(40101, errors.New("TestError")), 40101},
		{NewErrWithCode(40102, errors.New("TestError")), 40102},
	}
	for _, testdata := range tests {
		t.Run(testdata.err.Error(), func(t *testing.T) {
			resp := httptest.NewRecorder()

			Handler{logrus.New()}.Handle(resp, testdata.err)

			var expectedErr ErrorResponse
			for _, errResponse := range errResponses {
				if testdata.errCode == errResponse.Code {
					expectedErr = errResponse
					break
				}
			}
			expectedErr.Err = errorFormat.Error{testdata.err.Error()}
			if expectedErr.HTTPStatus != resp.Code {
				t.Error(test.Format("HTTP Status doesn't correspond.", strconv.Itoa(expectedErr.HTTPStatus), strconv.Itoa(resp.Code)))
			}

			var body ErrorResponse
			json.Unmarshal(resp.Body.Bytes(), &body)
			expectedErr.Time = body.Time
			if expectedErr.Code != body.Code {
				t.Error(test.Format("Body doesn't correspond to expected", strconv.Itoa(testdata.errCode), strconv.Itoa(body.Code)))
			}
		})
	}
}
