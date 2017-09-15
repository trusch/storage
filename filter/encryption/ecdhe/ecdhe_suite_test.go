package ecdhe_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEcdhe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ecdhe Suite")
}
