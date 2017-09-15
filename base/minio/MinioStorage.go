package minio

import (
	"errors"
	"io"
	"log"
	"time"

	minio "github.com/minio/minio-go"
)

// Storage implements the Storage interface with local files
type Storage struct {
	client *minio.Client
	bucket string
}

// NewStorage creates a new filestorage with the given base directory
// if baseDirectory doesn't exist, it is created like "mkdir -p $baseDirectory"
func NewStorage(url, bucket, accessKey, accessSecret string, ssl bool) (*Storage, error) {
	var cli *minio.Client
	for i := 0; i < 5; i++ {
		c, err := minio.New(url, accessKey, accessSecret, ssl)
		if err != nil {
			log.Print(err)
			time.Sleep(1 * time.Second)
			continue
		}
		cli = c
		break
	}

	location := "us-east-1"
	if exists, err := cli.BucketExists(bucket); err != nil || !exists {
		if err = cli.MakeBucket(bucket, location); err != nil {
			return nil, err
		}
	}

	return &Storage{cli, bucket}, nil
}

// GetReader returns a reader
func (fs *Storage) GetReader(id string) (io.ReadCloser, error) {
	return fs.client.GetObject(fs.bucket, id)
}

// GetWriter returns a writer
func (fs *Storage) GetWriter(id string) (io.WriteCloser, error) {
	return newMinioWriter(fs.client, fs.bucket, id), nil
}

// Has returns whether an entry exists
func (fs *Storage) Has(id string) bool {
	ob, err := fs.client.GetObject(fs.bucket, id)
	if err != nil {
		return false
	}
	_, err = ob.Stat()
	return err == nil
}

// Delete deletes an entry
func (fs *Storage) Delete(id string) error {
	if !fs.Has(id) {
		return errors.New("no such object")
	}
	return fs.client.RemoveObject(fs.bucket, id)
}

// List returns a list of all stored objects, limited by a prefix
func (fs *Storage) List(prefix string) ([]string, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	it := fs.client.ListObjects(fs.bucket, prefix, false, doneCh)
	res := make([]string, 0)
	for info := range it {
		res = append(res, info.Key)
	}
	return res, nil
}

func newMinioWriter(cli *minio.Client, bucket, id string) io.WriteCloser {
	r, w := io.Pipe()
	ch := make(chan struct{})
	res := &s3Writer{w, nil, ch}
	go func() {
		_, e := cli.PutObject(bucket, id, r, "application/binary")
		res.err = e
		close(ch)
	}()
	return res
}

type s3Writer struct {
	writer io.WriteCloser
	err    error
	done   chan struct{}
}

func (w *s3Writer) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *s3Writer) Close() error {
	err := w.writer.Close()
	if err != nil {
		return err
	}
	<-w.done
	return w.err
}
