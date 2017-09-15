package minio_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	uri    string = "127.0.0.1:9000"
	bucket string = "test"
	key    string = "YV4ZWPZ3HEUN4P21GUH0"
	secret string = "TU8KRW8dqstLv8t6tQHbAoOPqhPZWpyGzue46IK/"
)

func TestMinio(t *testing.T) {

	BeforeSuite(func() {
		dockerCall := fmt.Sprintf("docker run -p 9000:9000 --name minio --rm -d -e 'MINIO_ACCESS_KEY=%v' -e 'MINIO_SECRET_KEY=%v' minio/minio server /data", key, secret)
		Expect(exec.Command("bash", "-c", dockerCall).Run()).To(Succeed())
		exec.Command("bash", "-c", "while ! curl localhost:9000; do echo sleep; sleep 0.2; done; sleep 0.5;").Run()
	})

	AfterSuite(func() {
		Expect(exec.Command("bash", "-c", "docker stop minio").Run()).To(Succeed())
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Minio Suite")
}
