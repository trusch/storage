package aes_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	. "github.com/trusch/storage/base/file"
	. "github.com/trusch/storage/filter/encryption/aes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AES", func() {
	var (
		baseDirectory string = filepath.Join(os.TempDir(), "filestorage")
	)

	AfterEach(func() {
		Expect(os.RemoveAll(baseDirectory)).To(Succeed())
	})

	It("should be possible to save and load something", func() {
		storage := NewEncrypter(NewStorage(baseDirectory), "super-secret")
		writer, err := storage.GetWriter("test")
		Expect(err).NotTo(HaveOccurred())
		bs, err := writer.Write([]byte("foobar"))
		Expect(bs).To(Equal(6))
		Expect(err).NotTo(HaveOccurred())
		Expect(writer.Close()).To(Succeed())
		reader, err := storage.GetReader("test")
		Expect(err).NotTo(HaveOccurred())
		buf := &bytes.Buffer{}
		c, err := io.Copy(buf, reader)
		Expect(c).To(Equal(int64(6)))
		Expect(err).NotTo(HaveOccurred())
		Expect(buf.String()).To(Equal("foobar"))
		Expect(reader.Close()).To(Succeed())
	})

	It("should provide working has/delete methods", func() {
		storage := NewEncrypter(NewStorage(baseDirectory), "super-secret")
		Expect(storage.Has("test")).To(BeFalse())
		Expect(storage.Delete("test")).NotTo(Succeed())
		writer, err := storage.GetWriter("test")
		Expect(err).NotTo(HaveOccurred())
		bs, err := writer.Write([]byte("foobar"))
		Expect(bs).To(Equal(6))
		Expect(err).NotTo(HaveOccurred())
		Expect(writer.Close()).To(Succeed())
		Expect(storage.Has("test")).To(BeTrue())
		Expect(storage.Delete("test")).To(Succeed())
		Expect(storage.Has("test")).To(BeFalse())
	})
})
