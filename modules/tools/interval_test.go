package tools_test

import (
	"testing"

	"github.com/normegil/zookeeper-rest/modules/tools"
)

func TestNewIntegerInterval(t *testing.T) {
	tests := []struct {
		minIncluded bool
		min, max    int
		maxIncluded bool
	}{
		{true, 1, 2, true},
		{false, 1, 2, true},
		{true, 1, 2, false},
		{false, 1, 2, false},
		{true, 122594, 1082844, true},
		{false, 122594, 1082844, false},
		{true, 1, 1, true},
		{true, -1, 1, true},
		{true, -10, -5, true},
	}

	for _, test := range tests {
		i, err := tools.NewIntegerInterval(test.min, test.max, test.minIncluded, test.maxIncluded)
		if nil != err {
			t.Errorf("%+v: Error when running test: %s", test, err)
		} else {
			if i.Min() != test.min {
				t.Errorf("%+v: Min (%d) and test data doesn't correspond", i.Min(), test)
			}
			if i.Max() != test.max {
				t.Errorf("%+v: Max (%d) and test data doesn't correspond", i.Max(), test)
			}
			if i.MinIncluded() != test.minIncluded {
				t.Errorf("%+v: MinIncluded (%t) and test data doesn't correspond", i.MinIncluded(), test)
			}
			if i.MaxIncluded() != test.maxIncluded {
				t.Errorf("%+v: MaxIncluded (%t) and test data doesn't correspond", i.MaxIncluded(), test)
			}
		}
	}
}

func TestNewIntegerInterval_Errors(t *testing.T) {
	tests := []struct {
		minIncluded bool
		min, max    int
		maxIncluded bool
	}{
		{false, 1, 1, true},
		{true, 1, 1, false},
		{true, 2, 1, true},
		{true, -3, -10, true},
	}

	for _, test := range tests {
		_, err := tools.NewIntegerInterval(test.min, test.max, test.minIncluded, test.maxIncluded)
		if nil == err {
			t.Errorf("Error when running test with: %+v. Expected error.", test)
		}
	}
}

func TestParseIntervalInteger(t *testing.T) {
	tests := []struct {
		input       string
		minIncluded bool
		min, max    int
		maxIncluded bool
	}{
		{"[1;3]", true, 1, 3, true},
		{"[1;3[", true, 1, 3, false},
		{"]1;3]", false, 1, 3, true},
		{"]1;3[", false, 1, 3, false},
		{"]12;33[", false, 12, 33, false},
		{"]-3;-1[", false, -3, -1, false},
		{"]-25;-10[", false, -25, -10, false},
		{"]1042345;83762356[", false, 1042345, 83762356, false},
	}

	for _, test := range tests {
		i, err := tools.ParseIntervalInteger(test.input)
		if nil != err {
			t.Errorf("%+v: Error when running test: %s", test.input, err)
		} else {
			if i.Min() != test.min {
				t.Errorf("%+v: Min (%d) and test data doesn't correspond", i.Min(), test)
			}
			if i.Max() != test.max {
				t.Errorf("%+v: Max (%d) and test data doesn't correspond", i.Max(), test)
			}
			if i.MinIncluded() != test.minIncluded {
				t.Errorf("%+v: MinIncluded (%t) and test data doesn't correspond", i.MinIncluded(), test)
			}
			if i.MaxIncluded() != test.maxIncluded {
				t.Errorf("%+v: MaxIncluded (%t) and test data doesn't correspond", i.MaxIncluded(), test)
			}
		}
	}
}

func TestParseIntervalInteger_Errors(t *testing.T) {
	tests := []string{
		"",
		"abc",
		"[1,2]",
		"[1.2]",
		"(1;2]",
		"[1;2)",
		"[1;2(",
		")1;2]",
		"{1;2]",
		"[1;2}",
		"[1;2{",
		"}1;2]",
		"[a;2]",
		"[1;b]",
	}

	for _, test := range tests {
		_, err := tools.ParseIntervalInteger(test)
		if nil == err {
			t.Errorf("Error when running test with: %+v. Expected error.", test)
		}
	}
}

func TestInclude(t *testing.T) {
	tests := []struct {
		interval string
		value    int
		included bool
	}{
		{"[1;3]", 1, true},
		{"[1;3]", 2, true},
		{"[1;3]", 3, true},
		{"]1;2]", 1, false},
		{"]1;2]", 2, true},
		{"[1;2[", 1, true},
		{"[1;2[", 2, false},
		{"[1;2[", 10, false},
		{"[1;2[", -1, false},
		{"[-3;-1]", -10, false},
		{"[-3;-1]", -2, true},
	}
	for _, test := range tests {
		i, err := tools.ParseIntervalInteger(test.interval)
		if nil != err {
			t.Errorf("%+v: Error when running test: %s", test.interval, err)
		} else {
			included := i.Include(test.value)
			if included != test.included {
				t.Errorf("%s (%d): Included value (%t) is not equal to expected included value (%t)", test.interval, test.value, included, test.included)
			}
		}
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		interval string
		size     int
	}{
		{"[1;2]", 2},
		{"[1;2[", 1},
		{"[1;10]", 10},
		{"[1;10[", 9},
		{"[-5;5]", 11},
		{"[-5;-1]", 5},
	}
	for _, test := range tests {
		i, err := tools.ParseIntervalInteger(test.interval)
		if nil != err {
			t.Errorf("%+v: Error when running test: %s", test.interval, err)
		} else {
			size := i.Size()
			if size != test.size {
				t.Errorf("%s: Calculated size (%d) is not equal to expected size (%d)", test.interval, size, test.size)
			}
		}
	}
}

func TestLowestNumberIncluded(t *testing.T) {
	tests := []struct {
		interval string
		expected int
	}{
		{"[1;3]", 1},
		{"]1;3]", 2},
		{"]-1;3]", 0},
		{"[-5;3]", -5},
	}
	for _, test := range tests {
		i, err := tools.ParseIntervalInteger(test.interval)
		if nil != err {
			t.Errorf("%+v: Error when running test: %s", test.interval, err)
		} else {
			result := i.LowestNumberIncluded()
			if result != test.expected {
				t.Errorf("%s: Result for LowestNumberIncluded() (%d) is not equal to expected value (%d)", test.interval, result, test.expected)
			}
		}
	}
}

func TestHighestNumberIncluded(t *testing.T) {
	tests := []struct {
		interval string
		expected int
	}{
		{"[1;2]", 2},
		{"[1;2[", 1},
		{"[-10;-1]", -1},
	}
	for _, test := range tests {
		i, err := tools.ParseIntervalInteger(test.interval)
		if nil != err {
			t.Errorf("%+v: Error when running test: %s", test.interval, err)
		} else {
			result := i.HighestNumberIncluded()
			if result != test.expected {
				t.Errorf("%s: Result for HighestNumberIncluded() (%d) is not equal to expected value (%d)", test.interval, result, test.expected)
			}
		}
	}
}
