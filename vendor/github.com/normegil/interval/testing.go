package interval

import (
	"testing"

	"github.com/pkg/errors"
)

// Test_ParseIntervalInteger is a testing utility that will directly fail the test if the interval is not correct. Reserved for testing only.
func Test_ParseIntervalInteger(t testing.TB, s string) *IntervalInteger {
	interval, err := ParseIntervalInteger(s)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "Error while parsing %s", s))
	}
	return interval
}
