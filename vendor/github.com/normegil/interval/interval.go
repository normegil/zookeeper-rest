// Interval implementations for different type of intervals, defined by upper & lower bound, wichi can be exluded or included in the interval.
// Intervals are absolutes and cannot be negative, meaning that [2;1] will not be supported.
package interval

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// IntervalInteger represents an interval of integers
type IntervalInteger struct {
	min         int
	max         int
	minIncluded bool
	maxIncluded bool
}

// NewIntegerInterval should be used to create a new interval of integers
func NewIntegerInterval(min, max int, minIncluded, maxIncluded bool) (*IntervalInteger, error) {
	if min == max && (!minIncluded || !maxIncluded) {
		return nil, fmt.Errorf("Interval could not be built: single element %d is excluded", min)
	}
	if max < min {
		return nil, fmt.Errorf("Max %d is inferior to Min %d", max, min)
	}

	return &IntervalInteger{
		min:         min,
		max:         max,
		minIncluded: minIncluded,
		maxIncluded: maxIncluded,
	}, nil
}

/*
ParseIntervalInteger can be used to create a interval of integers from a string.
Supported format are:
	[1;2]
	]1;2[
	[1;2[
	]1;2]
*/
func ParseIntervalInteger(s string) (*IntervalInteger, error) {
	externalSign := `(\[|\])`
	separator := `;`
	number := `-?[0-9]*`
	intervalRegex := regexp.MustCompile(`^` + externalSign + number + separator + number + externalSign + `$`)

	if !intervalRegex.MatchString(s) {
		return nil, fmt.Errorf("Not a valid interval: %s", s)
	}

	var minIncluded bool
	firstCharacter := string(s[0])
	if "[" == firstCharacter {
		minIncluded = true
	}

	lastCharacter := string(s[len(s)-1])
	var maxIncluded bool
	if "]" == lastCharacter {
		maxIncluded = true
	}

	numbersAndSeparator := s[1 : len(s)-1]
	numbers := strings.Split(numbersAndSeparator, ";")
	min, err := strconv.Atoi(numbers[0])
	if nil != err {
		return nil, errors.Wrapf(err, "Could not get minimum value from: %s", numbers[0])
	}

	max, err := strconv.Atoi(numbers[1])
	if nil != err {
		return nil, errors.Wrapf(err, "Could not get maximum value from: %s", numbers[1])
	}

	return NewIntegerInterval(min, max, minIncluded, maxIncluded)
}

// Include test if the interval include the given integer
func (i IntervalInteger) Include(toTest int) bool {
	if toTest > i.Min() && toTest < i.Max() {
		return true
	}
	if toTest == i.Min() && i.MinIncluded() {
		return true
	}
	if toTest == i.Max() && i.MaxIncluded() {
		return true
	}
	return false
}

// Min return the lower bound of the interval
func (i IntervalInteger) Min() int {
	return i.min
}

// Max return the higher bound of the interval
func (i IntervalInteger) Max() int {
	return i.max
}

// MinIncluded return if the minimum should be included in the interval
func (i IntervalInteger) MinIncluded() bool {
	return i.minIncluded
}

// MinIncluded return if the minimum should be included in the interval
func (i IntervalInteger) MaxIncluded() bool {
	return i.maxIncluded
}

// LowestNumberIncluded will return the lowest number included, the lower bound if it's included or the lower bound + 1 if it's not
func (i IntervalInteger) LowestNumberIncluded() int {
	if i.MinIncluded() {
		return i.Min()
	}
	return i.Min() + 1
}

// HighestNumberIncluded will return the highest number included, the higher bound if it's included or the higher bound - 1 if it's not
func (i IntervalInteger) HighestNumberIncluded() int {
	if i.MaxIncluded() {
		return i.Max()
	}

	return i.Max() - 1
}

// Size return the number of integers included in the interval
func (i IntervalInteger) Size() int {
	return i.HighestNumberIncluded() - i.LowestNumberIncluded() + 1
}
