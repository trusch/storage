package gcs_test

import (
	"bytes"
	"context"
	"io"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"

	. "github.com/trusch/storage/base/gcs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoogleBucketStorage", func() {
	var (
		projectID = "webvariants-playground"
		bucketID  = "wv-backup-test"
	)

	AfterEach(func() {
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		Expect(err).NotTo(HaveOccurred())
		bucket := client.Bucket(bucketID)
		it := bucket.Objects(ctx, nil)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				Fail(err.Error())
			}
			Expect(bucket.Object(attrs.Name).Delete(ctx)).To(Succeed())
		}
		Expect(bucket.Delete(ctx)).To(Succeed())
		Expect(client.Close()).To(Succeed())
	})

	It("should be possible to save and load something", func() {
		store, err := NewStorage(projectID, bucketID)
		Expect(err).NotTo(HaveOccurred())
		writer, err := store.GetWriter("test")
		Expect(err).NotTo(HaveOccurred())
		bs, err := writer.Write([]byte("foobar"))
		Expect(bs).To(Equal(6))
		Expect(err).NotTo(HaveOccurred())
		Expect(writer.Close()).To(Succeed())
		reader, err := store.GetReader("test")
		Expect(err).NotTo(HaveOccurred())
		buf := &bytes.Buffer{}
		c, err := io.Copy(buf, reader)
		Expect(c).To(Equal(int64(6)))
		Expect(err).NotTo(HaveOccurred())
		Expect(buf.String()).To(Equal("foobar"))
	})

	It("should provide working has/delete methods", func() {
		store, err := NewStorage(projectID, bucketID)
		Expect(err).NotTo(HaveOccurred())
		Expect(store.Has("test")).To(BeFalse())
		Expect(store.Delete("test")).NotTo(Succeed())
		writer, err := store.GetWriter("test")
		Expect(err).NotTo(HaveOccurred())
		bs, err := writer.Write([]byte("foobar"))
		Expect(bs).To(Equal(6))
		Expect(err).NotTo(HaveOccurred())
		Expect(writer.Close()).To(Succeed())
		Expect(store.Has("test")).To(BeTrue())
		Expect(store.Delete("test")).To(Succeed())
		Expect(store.Has("test")).To(BeFalse())
	})

	It("should provide a list method", func() {
		store, err := NewStorage(projectID, bucketID)
		Expect(err).NotTo(HaveOccurred())
		writer, err := store.GetWriter("a")
		Expect(err).NotTo(HaveOccurred())
		Expect(writer.Close()).To(Succeed())
		writer, err = store.GetWriter("b")
		Expect(err).NotTo(HaveOccurred())
		Expect(writer.Close()).To(Succeed())
		writer, err = store.GetWriter("bb")
		Expect(err).NotTo(HaveOccurred())
		Expect(writer.Close()).To(Succeed())

		objects, err := store.List("")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(objects)).To(Equal(3))
		Expect(objects[0]).To(Equal("a"))
		Expect(objects[1]).To(Equal("b"))
		Expect(objects[2]).To(Equal("bb"))

		objects, err = store.List("b")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(objects)).To(Equal(2))
		Expect(objects[0]).To(Equal("b"))
		Expect(objects[1]).To(Equal("bb"))
	})
})
