package tools

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type IntervalInteger struct {
	min         int
	max         int
	minIncluded bool
	maxIncluded bool
}

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

func (i IntervalInteger) Min() int {
	return i.min
}

func (i IntervalInteger) Max() int {
	return i.max
}

func (i IntervalInteger) MinIncluded() bool {
	return i.minIncluded
}

func (i IntervalInteger) MaxIncluded() bool {
	return i.maxIncluded
}

func (i IntervalInteger) LowestNumberIncluded() int {
	if i.MinIncluded() {
		return i.Min()
	}
	return i.Min() + 1
}

func (i IntervalInteger) HighestNumberIncluded() int {
	if i.MaxIncluded() {
		return i.Max()
	}

	return i.Max() - 1
}

func (i IntervalInteger) Size() int {
	return i.HighestNumberIncluded() - i.LowestNumberIncluded() + 1
}
