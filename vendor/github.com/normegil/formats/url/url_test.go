package url_test

import (
	"encoding/json"
	"net/url"
	"testing"

	formatPkg "github.com/normegil/formats/url"
	"github.com/pkg/errors"
)

func TestMarshal(t *testing.T) {
	testcases := []struct {
		Input  string
		Output string
	}{
		{"http://www.example.com", "\"http://www.example.com\""},
		{"https://www.example.com", "\"https://www.example.com\""},
		{"http://www.example.com/test", "\"http://www.example.com/test\""},
		{"http://www.example.com/test/t", "\"http://www.example.com/test/t\""},
		{"http://www.example.com/test?key=value", "\"http://www.example.com/test?key=value\""},
		{"http://www.example.com/test?key=value&key2=value2", "\"http://www.example.com/test?key=value\\u0026key2=value2\""},
	}
	for _, testdata := range testcases {
		t.Run(testdata.Input, func(t *testing.T) {
			parsed, err := url.Parse(testdata.Input)
			if err != nil {
				t.Fatal(errors.Wrapf(err, "Parsing %s", testdata.Input))
			}
			toFormat := formatPkg.URL{parsed}
			jsonURL, err := json.Marshal(toFormat)
			if err != nil {
				t.Fatal(errors.Wrapf(err, "Marshal URL %+v", toFormat))
			}

			if testdata.Output != string(jsonURL) {
				t.Errorf("Expected (%s) is not equals to computed (%s)", testdata.Output, string(jsonURL))
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	testcases := []string{
		"http://www.example.com",
		"https://www.example.com",
		"http://www.example.com/test",
		"http://www.example.com/test/t",
		"http://www.example.com/test?key=value",
		"http://www.example.com/test?key=value&key2=value2",
		"http://www.example.com/test?key=value\u0026key2=value2",
	}
	for _, testdata := range testcases {
		t.Run(testdata, func(t *testing.T) {
			toParse := "\"" + testdata + "\""
			var parsedURL formatPkg.URL
			if err := json.Unmarshal([]byte(toParse), &parsedURL); nil != err {
				t.Fatal(errors.Wrapf(err, "Parsing %s", testdata))
			}

			if testdata != parsedURL.String() {
				t.Error("Expected (%s) is not equals to computed (%s)", testdata, parsedURL)
			}
		})
	}
}
