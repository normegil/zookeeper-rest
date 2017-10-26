// Small library containing utilities to use with streams in Go
package streamutils

import (
	"bytes"
	"io"
)

// ToBytes convert an io.Reader into a byte array
func ToBytes(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
