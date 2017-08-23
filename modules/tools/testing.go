package tools

import (
	"testing"

	"github.com/pkg/errors"
)

func Test_ParseIntervalInteger(t testing.TB, s string) *IntervalInteger {
	interval, err := ParseIntervalInteger(s)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "Error while parsing %s", s))
	}
	return interval
}
