package minio_test

import (
	"bytes"
	"io"

	. "github.com/trusch/storage/base/minio"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MinioStorage", func() {

	It("should be possible to save and load something", func() {
		storage, err := NewStorage(uri, bucket, key, secret, false)
		Expect(err).NotTo(HaveOccurred())
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
		// cleanup
		Expect(storage.Delete("test")).To(Succeed())
	})

	It("should provide working has/delete methods", func() {
		storage, err := NewStorage(uri, bucket, key, secret, false)
		Expect(err).NotTo(HaveOccurred())
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

	It("should provide a list method", func() {
		storage, err := NewStorage(uri, bucket, key, secret, false)
		Expect(err).NotTo(HaveOccurred())
		writer, err := storage.GetWriter("a")
		Expect(err).NotTo(HaveOccurred())
		writer.Write([]byte("foobar"))
		Expect(writer.Close()).To(Succeed())
		writer, err = storage.GetWriter("b")
		Expect(err).NotTo(HaveOccurred())
		writer.Write([]byte("foobar"))
		Expect(writer.Close()).To(Succeed())
		writer, err = storage.GetWriter("bb")
		Expect(err).NotTo(HaveOccurred())
		writer.Write([]byte("foobar"))
		Expect(writer.Close()).To(Succeed())

		objects, err := storage.List("")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(objects)).To(Equal(3))
		Expect(objects[0]).To(Equal("a"))
		Expect(objects[1]).To(Equal("b"))
		Expect(objects[2]).To(Equal("bb"))

		objects, err = storage.List("b")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(objects)).To(Equal(2))
		Expect(objects[0]).To(Equal("b"))
		Expect(objects[1]).To(Equal("bb"))

		// cleanup
		Expect(storage.Delete("a")).To(Succeed())
		Expect(storage.Delete("b")).To(Succeed())
		Expect(storage.Delete("bb")).To(Succeed())
	})
})
