package streamutils_test

import (
	"io"
	"strings"
	"testing"

	"github.com/normegil/streamutils"
)

func TestToByte(t *testing.T) {
	tests := []struct {
		input  io.Reader
		output string
	}{
		{strings.NewReader("test"), "test"},
	}

	for _, test := range tests {
		result := streamutils.ToBytes(test.input)
		convertedResult := string(result)
		if test.output != convertedResult {
			t.Errorf("ToByte should use stream and translate it to bytes {Result:%+v;Expected:%+v}", result, test.output)
		}
	}
}
