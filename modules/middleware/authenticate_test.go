package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/normegil/resterrors"
	"github.com/normegil/zookeeper-rest/modules/middleware"
	"github.com/normegil/zookeeper-rest/modules/model"
	"github.com/normegil/zookeeper-rest/modules/model/dao"
	"github.com/normegil/zookeeper-rest/modules/test"
	"github.com/sirupsen/logrus"
)

type testHandler struct {
	test *testing.T
}

func (t testHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	user := request.Context().Value(middleware.USER_CTX_KEY)
	toWrite, err := json.Marshal(user)
	if err != nil {
		t.test.Fatalf("Could not Marshal user: %+v", user)
	}
	response.Write(toWrite)
}

func TestRequestAuthenticator(t *testing.T) {
	testcases := []struct {
		toAuthenticate model.User
		dbUsers        []model.User
		expectedCode   int
	}{
		{
			toAuthenticate: model.UserImpl{"test", "test"},
			expectedCode:   200,
			dbUsers: []model.User{
				model.UserImpl{"test", "test"},
			},
		},
		{
			toAuthenticate: model.UserImpl{"", "test"},
			expectedCode:   40103,
			dbUsers:        []model.User{},
		},
	}
	for _, testdata := range testcases {
		t.Run("", func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://localhost/", strings.NewReader(""))
			request.SetBasicAuth(testdata.toAuthenticate.Name(), testdata.toAuthenticate.Password())
			result := httptest.NewRecorder()

			handler := middleware.RequestAuthenticator(logrus.New(), &dao.Test_UserDAO{testdata.dbUsers}, testHandler{t})
			handler.ServeHTTP(result, request)

			if 200 == result.Code {
				if testdata.expectedCode != result.Code {
					t.Fatalf(test.Format("Code doesn't correspond to expectedCode", testdata.expectedCode, result.Code))
				}
				usr := &model.UserImpl{}
				err := json.Unmarshal(result.Body.Bytes(), usr)
				if err != nil {
					t.Fatal(err)
				}
				if testdata.toAuthenticate.Name() != usr.Name() || testdata.toAuthenticate.Password() != usr.Password() {
					t.Error(test.Format("Loaded user doesn't correspond to requested user", testdata.toAuthenticate.String(), usr.String()))
				}
			} else {
				resp := &resterrors.ErrorResponse{}
				err := json.Unmarshal(result.Body.Bytes(), resp)
				if err != nil {
					t.Fatal(err)
				}
				if testdata.expectedCode != resp.Code {
					t.Fatalf(test.Format("Code doesn't correspond to expectedCode", testdata.expectedCode, resp.Code))
				}
			}
		})
	}
}
