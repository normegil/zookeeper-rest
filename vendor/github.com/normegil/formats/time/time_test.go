package time_test

import (
	"testing"
	"time"
)

func TestTimeToString(t *testing.T) {
	cases := []time.Time{
		time.Now(),
		time.Now().Add(48 * time.Hour),
	}
	for _, testTime := range cases {
		value := time.Time(testTime).String()
		expected := testTime.Format(time.RFC3339)
		if expected != value {
			t.Errorf("Returned message (%s) don't correspond to expected message (%s)", expected, value)
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
		value, err := time.Time(params.input).MarshalJSON()
		if nil != err {
			t.Fatal(err.Error())
		}
		output := "\"" + params.input.Format(time.RFC3339) + "\""
		if output != string(value) {
			t.Errorf("Malformed JSON\n\tExpected:%s\n\tGot:%s", output, string(value))
		}
	}
}
