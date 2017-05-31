package tools_test

import (
	. "github.com/normegil/zookeeper-rest/modules/tools"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("ToBytes", func() {
	Context("When passed an io.Reader", func() {
		It("should return the full content of the reader as bytes", func() {
			text := "Test\nTest2"
			stream := strings.NewReader(text)
			Expect(ToBytes(stream)).To(Equal([]byte(text)))
		})
	})
})
