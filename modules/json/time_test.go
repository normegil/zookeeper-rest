package json_test

import (
	"testing"
	"time"

	"github.com/normegil/zookeeper-rest/modules/json"
	"github.com/normegil/zookeeper-rest/modules/test"
)

func TestTimeToString(t *testing.T) {
	cases := []time.Time{
		time.Now(),
		time.Now().Add(48 * time.Hour),
	}
	for _, testTime := range cases {
		value := json.JSONTime(testTime).String()
		expected := testTime.Format(time.RFC3339)
		if expected != value {
			t.Error(test.Format("JSONTime.String()", "Returned message don't correspond to expectd output", expected, value))
		}
	}
}

func TestTimeMarshallJSON(t *testing.T) {
	cases := []struct {
		testName string
		input    time.Time
	}{
		{"JSON - Classic case", time.Now()},
	}
	for _, params := range cases {
		value, err := json.JSONTime(params.input).MarshalJSON()
		if nil != err {
			t.Fatal(err.Error())
		}
		output := "\"" + params.input.Format(time.RFC3339) + "\""
		if output != string(value) {
			t.Error(test.Format(params.testName, "Malformed JSON", output, string(value)))
		}
	}
}
