package formats_test

import (
	"encoding/json"
	"testing"

	"github.com/normegil/zookeeper-rest/modules/formats"
	"github.com/normegil/zookeeper-rest/modules/test"
)

func TestErrorToString(t *testing.T) {
	cases := []string{
		"",
		"Test",
	}
	for _, errMsg := range cases {
		err := formats.Error{errMsg}
		if errMsg != err.Error() {
			t.Error(test.Format("Returned message don't correspond", errMsg, err.Error()))
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
			bytes, err := formats.Error{params.input}.MarshalJSON()
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
				t.Error(test.Format("Malformed JSON", string(expectedBytes), string(bytes)))
			}
		})
	}
}
