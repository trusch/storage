package ecdhe_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	. "github.com/trusch/storage/base/file"
	. "github.com/trusch/storage/filter/encryption/ecdhe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ECDHE", func() {
	var (
		baseDirectory string = filepath.Join(os.TempDir(), "filestorage")
		privKey              = []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIEAfatPllysuYKwdXUDV7lykSYJwm7w172rBIZK8efJXoAoGCCqGSM49
AwEHoUQDQgAEc9ba+mhTJdL9rnIkFnPcgLpSlFJ/2tFJ38+PSOQ/1OA2zEYTRXT4
h5LlKFNB1dGWYKhKsXkiocY4POI1Y59rHQ==
-----END EC PRIVATE KEY-----
`)
		pubKey = []byte(`-----BEGIN CERTIFICATE-----
MIIBZzCCAQ2gAwIBAgIBATAKBggqhkjOPQQDAjAiMRAwDgYDVQQKEwdBY21lIENv
MQ4wDAYDVQQDEwUuL3BraTAeFw0xNzA5MDUwNTQ5MjJaFw0yNzA5MDMwNTQ5MjJa
MCExEDAOBgNVBAoTB0FjbWUgQ28xDTALBgNVBAMTBHRlc3QwWTATBgcqhkjOPQIB
BggqhkjOPQMBBwNCAARz1tr6aFMl0v2uciQWc9yAulKUUn/a0Unfz49I5D/U4DbM
RhNFdPiHkuUoU0HV0ZZgqEqxeSKhxjg84jVjn2sdozUwMzAOBgNVHQ8BAf8EBAMC
BaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAKBggqhkjOPQQD
AgNIADBFAiAunBCPeVQtN2Ytsml3HJ8OztlLFqPwxu4km+aVx+LlhAIhAJpvQwvA
UKY+VRlEoIH1Y5DcCvGMvByMOTNBcGsvgm9M
-----END CERTIFICATE-----
`)
	)

	AfterEach(func() {
		Expect(os.RemoveAll(baseDirectory)).To(Succeed())
	})

	It("should be possible to save and load something", func() {
		storage, err := NewEncrypter(NewStorage(baseDirectory), pubKey, privKey)
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
		Expect(reader.Close()).To(Succeed())
	})

	It("should provide working has/delete methods", func() {
		storage, err := NewEncrypter(NewStorage(baseDirectory), pubKey, privKey)
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
})
