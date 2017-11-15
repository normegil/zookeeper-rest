package error_test

import (
	"encoding/json"
	"testing"

	"github.com/normegil/formats/error"
)

func TestErrorToString(t *testing.T) {
	cases := []string{
		"",
		"Test",
	}
	for _, errMsg := range cases {
		err := error.Error{errMsg}
		if errMsg != err.Error() {
			t.Errorf("Returned message (%s) don't correspond to expected error message (%s)", errMsg, err.Error())
		}
	}
}

func TestErrorMarshallJSON(t *testing.T) {
	cases := []struct {
		testName string
		input    string
	}{
		{"JSON - Empty field", ""},
		{"JSON - Classic case", "Test"},
	}
	for _, params := range cases {
		t.Run(params.testName, func(t *testing.T) {
			bytes, err := error.Error{params.input}.MarshalJSON()
			if nil != err {
				t.Fatal(err.Error())
			}

			expected := make(map[string]interface{})
			expected["@type"] = "BaseError"
			expected["message"] = params.input
			expectedBytes, err := json.Marshal(expected)
			if err != nil {
				panic(err)
			}

			if string(expectedBytes) != string(bytes) {
				t.Error("Malformed JSON\n\tExpected:%s\n\tGot:%s", string(expectedBytes), string(bytes))
			}
		})
	}
}
