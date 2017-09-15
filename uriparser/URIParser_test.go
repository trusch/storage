package uriparser_test

import (
	"bytes"
	"io"
	"os"

	"github.com/trusch/storage"
	. "github.com/trusch/storage/uriparser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testStorage(storage storage.Storage) {
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
}

var _ = Describe("URIParser", func() {
	AfterEach(func() {
		Expect(os.RemoveAll("/tmp/backups")).To(Succeed())
	})
	It("should parse simple file:// scheme", func() {
		storage, err := NewFromURI("file:///tmp/backups", nil)
		Expect(err).NotTo(HaveOccurred())
		testStorage(storage)
	})
	It("should parse gzip+file:// scheme", func() {
		storage, err := NewFromURI("gzip+file:///tmp/backups", nil)
		Expect(err).NotTo(HaveOccurred())
		testStorage(storage)
	})
	It("should parse snappy+file:// scheme", func() {
		storage, err := NewFromURI("snappy+file:///tmp/backups", nil)
		Expect(err).NotTo(HaveOccurred())
		testStorage(storage)
	})
	It("should parse xz+file:// scheme", func() {
		storage, err := NewFromURI("xz+file:///tmp/backups", nil)
		Expect(err).NotTo(HaveOccurred())
		testStorage(storage)
	})
	It("should parse gzip+aes+file:// scheme", func() {
		storage, err := NewFromURI("gzip+aes+file:///tmp/backups", Options{"key": "my-aes-key"})
		Expect(err).NotTo(HaveOccurred())
		testStorage(storage)
	})
	It("should parse snappy+aes+file:// scheme", func() {
		storage, err := NewFromURI("snappy+aes+file:///tmp/backups", Options{"key": "my-aes-key"})
		Expect(err).NotTo(HaveOccurred())
		testStorage(storage)
	})
	It("should parse xz+aes+file:// scheme", func() {
		storage, err := NewFromURI("xz+aes+file:///tmp/backups", Options{"key": "my-aes-key"})
		Expect(err).NotTo(HaveOccurred())
		testStorage(storage)
	})
})
