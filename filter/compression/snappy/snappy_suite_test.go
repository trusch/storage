package snappy_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSnappy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Snappy Suite")
}
