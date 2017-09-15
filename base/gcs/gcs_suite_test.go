package gcs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGcs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gcs Suite")
}
