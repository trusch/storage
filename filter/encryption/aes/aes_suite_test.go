package aes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aes Suite")
}
