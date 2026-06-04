package koito

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKoito(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Koito Suite")
}
