package json_test

import (
	"errors"
	"testing"

	"github.com/normegil/aphrodite/modules/json"
	"github.com/normegil/aphrodite/modules/test"
)

func TestErrorToString(t *testing.T) {
	cases := []string{
		"",
		"Test",
	}
	for _, errMsg := range cases {
		jsonErr := json.ErrorJSON{errors.New(errMsg)}
		if errMsg != jsonErr.Error() {
			t.Error(test.Format("ErrorJSON.Error()", "Returned message don't correspond", errMsg, jsonErr.Error()))
		}
	}
}

func TestErrorMarshallJSON(t *testing.T) {
	cases := []struct {
		testName string
		input    string
		output   string
	}{
		{"JSON - Empty field", "", "\"\""},
		{"JSON - Classic case", "Test", "\"Test\""},
	}
	for _, params := range cases {
		bytes, err := json.ErrorJSON{errors.New(params.input)}.MarshalJSON()
		if nil != err {
			t.Fatal(err.Error())
		}
		if params.output != string(bytes) {
			t.Error(test.Format(params.testName, "Malformed JSON", params.output, params.input))
		}
	}
}
