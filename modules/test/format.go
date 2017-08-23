package test

import "fmt"

type PrefixedStringer interface {
	PrefixedString(prefix string) string
}

func Format(msg string, expected, got interface{}) string {
	if "" == msg {
		msg = "Values should be equals"
	}
	expectedStr := toString(expected, "\t\t")
	gotStr := toString(got, "\t\t")
	return fmt.Sprintf("%s\n\tExpected:\n%s\n\tGot:\n%s", msg, expectedStr, gotStr)
}

func toString(object interface{}, prefix string) string {
	if nil == object {
		return "nil"
	}

	if objPrefixStringer, ok := object.(PrefixedStringer); ok {
		return objPrefixStringer.PrefixedString(prefix)
	} else if objString, ok := object.(fmt.Stringer); ok {
		return objString.String()
	}
	return fmt.Sprintf("%+v", object)
}
