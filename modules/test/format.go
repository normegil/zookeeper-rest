package test

import "fmt"

func Format(testName, msg, expected, got string) string {
	return fmt.Sprintf("%s\n\t%s\n\t\tExpected: %+v\n\t\tGot: %+v", testName, msg, expected, got)
}
