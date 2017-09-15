package mongo_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMongo(t *testing.T) {
	BeforeSuite(func() {
		Expect(exec.Command("bash", "-c", "docker run -p 27017:27017 --name mongodb --rm -d mongo").Run()).To(Succeed())
	})

	AfterSuite(func() {
		Expect(exec.Command("bash", "-c", "docker stop mongodb").Run()).To(Succeed())
	})
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongo Suite")
}
