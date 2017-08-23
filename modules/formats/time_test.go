package formats_test

import (
	"testing"
	"time"

	"github.com/normegil/zookeeper-rest/modules/formats"
	"github.com/normegil/zookeeper-rest/modules/test"
)

func TestTimeToString(t *testing.T) {
	cases := []time.Time{
		time.Now(),
		time.Now().Add(48 * time.Hour),
	}
	for _, testTime := range cases {
		value := formats.Time(testTime).String()
		expected := testTime.Format(time.RFC3339)
		if expected != value {
			t.Error(test.Format("Returned message don't correspond to expectd output", expected, value))
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
		value, err := formats.Time(params.input).MarshalJSON()
		if nil != err {
			t.Fatal(err.Error())
		}
		output := "\"" + params.input.Format(time.RFC3339) + "\""
		if output != string(value) {
			t.Error(test.Format("Malformed JSON", output, string(value)))
		}
	}
}
