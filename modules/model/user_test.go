package model_test

import (
	"testing"

	"github.com/normegil/zookeeper-rest/modules/model"
	"github.com/normegil/zookeeper-rest/modules/test"
)

func TestCheckPassword(t *testing.T) {
	testcases := []struct {
		Pass           string
		ToCheck        string
		ExpectedResult bool
	}{
		{"a", "a", true},
		{"b", "a", false},
		{"@$^µ", "@$^µ", true},
		{"@$^µ", "a", false},
	}
	for _, testdata := range testcases {
		t.Run(testdata.Pass+":"+testdata.ToCheck, func(t *testing.T) {
			result := model.UserImpl{Pass: testdata.Pass}.Check(testdata.ToCheck)
			if testdata.ExpectedResult != result {
				t.Error(test.Format("Results for password checking operation don't match [Pass:'"+testdata.Pass+"';Checked:'"+testdata.ToCheck+"']", testdata.ExpectedResult, result))
			}
		})
	}
}
