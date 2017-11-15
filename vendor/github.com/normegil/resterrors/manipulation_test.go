package resterrors_test

import (
	"testing"

	"github.com/normegil/resterrors"
	"github.com/pkg/errors"
)

func TestSearchThroughCauses(t *testing.T) {
	testErr := errors.New("Test")
	errCode := resterrors.NewErrWithCode(100, testErr)
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
				_, ok := e.(resterrors.ErrWithCode)
				return ok
			},
		},
		{
			Input:    nil,
			Expected: nil,
			SearchFunction: func(e error) bool {
				_, ok := e.(resterrors.ErrWithCode)
				return ok
			},
		},
		{
			Input:    errors.Wrapf(testErr, ""),
			Expected: nil,
			SearchFunction: func(e error) bool {
				_, ok := e.(resterrors.ErrWithCode)
				return ok
			},
		},
	}
	for _, testdata := range testcases {
		t.Run("", func(t *testing.T) {
			found := resterrors.SearchThroughCauses(testdata.Input, testdata.SearchFunction)
			if testdata.Expected != found {
				t.Error("Searched error (%+v) was not the found error (%+v)", found, testdata.Expected)
			}
		})
	}
}
