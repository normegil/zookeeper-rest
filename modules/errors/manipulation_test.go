package errors_test

import (
	"testing"

	errPkg "github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/normegil/zookeeper-rest/modules/test"
	"github.com/pkg/errors"
)

func TestSearchThroughCauses(t *testing.T) {
	testErr := errors.New("Test")
	errCode := errPkg.NewErrWithCode(100, testErr)
	testcases := []struct {
		Input          error
		Expected       error
		SearchFunction func(error) bool
	}{
		{
			Input:    testErr,
			Expected: testErr,
			SearchFunction: func(e error) bool {
				return true
			},
		},
		{
			Input:    errors.Wrapf(errors.Wrapf(errCode, ""), ""),
			Expected: errCode,
			SearchFunction: func(e error) bool {
				_, ok := e.(errPkg.ErrWithCode)
				return ok
			},
		},
		{
			Input:    nil,
			Expected: nil,
			SearchFunction: func(e error) bool {
				_, ok := e.(errPkg.ErrWithCode)
				return ok
			},
		},
		{
			Input:    errors.Wrapf(testErr, ""),
			Expected: nil,
			SearchFunction: func(e error) bool {
				_, ok := e.(errPkg.ErrWithCode)
				return ok
			},
		},
	}
	for _, testdata := range testcases {
		t.Run("", func(t *testing.T) {
			found := errPkg.SearchThroughCauses(testdata.Input, testdata.SearchFunction)
			if testdata.Expected != found {
				t.Error(test.Format("Searched error was not the found error", testdata.Expected, found))
			}
		})
	}
}
