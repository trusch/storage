package xz_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestXz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Xz Suite")
}
