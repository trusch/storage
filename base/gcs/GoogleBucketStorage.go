package gcs

import (
	"context"
	"io"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

// Storage implements the Storage interface with local files
type Storage struct {
	ctx    context.Context
	bucket *storage.BucketHandle
}

// NewStorage creates a new filestorage with the given base directory
// if baseDirectory doesn't exist, it is created like "mkdir -p $baseDirectory"
func NewStorage(projectID, bucketID string) (*Storage, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bucket := client.Bucket(bucketID)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		if err = bucket.Create(ctx, projectID, nil); err != nil {
			return nil, err
		}
	}
	return &Storage{ctx, bucket}, nil
}

// GetReader returns a reader
func (fs *Storage) GetReader(id string) (io.ReadCloser, error) {
	return fs.bucket.Object(id).NewReader(fs.ctx)
}

// GetWriter returns a writer
func (fs *Storage) GetWriter(id string) (io.WriteCloser, error) {
	return fs.bucket.Object(id).NewWriter(fs.ctx), nil
}

// Has returns whether an entry exists
func (fs *Storage) Has(id string) bool {
	_, err := fs.bucket.Object(id).Attrs(fs.ctx)
	return err == nil
}

// Delete deletes an entry
func (fs *Storage) Delete(id string) error {
	return fs.bucket.Object(id).Delete(fs.ctx)
}

// List lists all stored objects with the given prefix
func (fs *Storage) List(prefix string) ([]string, error) {
	it := fs.bucket.Objects(fs.ctx, &storage.Query{Prefix: prefix})
	res := make([]string, 0)
	for {
		info, err := it.Next()
		if err == nil {
			res = append(res, info.Name)
		} else {
			if err == iterator.Done {
				break
			}
			return nil, err
		}
	}
	return res, nil
}
