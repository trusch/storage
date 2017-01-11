package file

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/trusch/storage/common"
)

// Storage creates the apropriate store from an URI
type Storage struct {
	base string
}

// NewStorage creates a new storage from a URI
func NewStorage(base string) (*Storage, error) {
	err := os.MkdirAll(base, 0700)
	if err != nil {
		return nil, err
	}
	return &Storage{base}, nil
}

// Put saves a byteslice to the db.
// Example: Save("/foo/bar", []byte{1,2,3})
func (store *Storage) Put(bucket, key string, value []byte) error {
	path := filepath.Join(store.base, bucket, key)
	return ioutil.WriteFile(path, value, 0600)
}

// Get loads data from a key
func (store *Storage) Get(bucket, key string) ([]byte, error) {
	path := filepath.Join(store.base, bucket, key)
	return ioutil.ReadFile(path)
}

// Delete deletes a value from the db
func (store *Storage) Delete(bucket, key string) error {
	path := filepath.Join(store.base, bucket, key)
	if _, err := os.Stat(path); err != nil {
		if _, err = os.Stat(filepath.Join(store.base, bucket)); err != nil {
			return common.Error(common.BucketNotFound, err)
		}
		return nil
	}
	return os.Remove(path)
}

// CreateBucket creates a bucket
func (store *Storage) CreateBucket(bucket string) error {
	path := filepath.Join(store.base, bucket)
	return os.MkdirAll(path, 0700)
}

// DeleteBucket deletes a bucket
func (store *Storage) DeleteBucket(bucket string) error {
	path := filepath.Join(store.base, bucket)
	if _, err := os.Stat(path); err != nil {
		return common.Error(common.BucketNotFound, err)
	}
	return os.RemoveAll(path)
}

// List returns all Entries of a directory
// optionally provide arguments to specifiy a key offset and a key limit
// Example: List("/foo", "abc", "xyz") -> DocInfo{Key: abc} ... DocInfo{Key: ggg} ... DocInfo{Key: xyz}
func (store *Storage) List(bucket string, opts *common.ListOpts) (chan *common.DocInfo, error) {
	path := filepath.Join(store.base, bucket)
	if _, err := os.Stat(path); err != nil {
		return nil, common.Error(common.BucketNotFound, err)
	}
	if opts == nil {
		opts = &common.ListOpts{}
	}
	ch := make(chan *common.DocInfo, 64)

	go func() {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			key := filepath.Base(path)
			val, err := ioutil.ReadFile(path)
			if err != nil {
				log.Print(err)
				return nil
			}
			switch {
			case opts.Prefix != "":
				{
					if strings.HasPrefix(key, opts.Prefix) {
						ch <- &common.DocInfo{Key: key, Value: val}
					}
				}
			case opts.Start != "":
				{
					if strings.Compare(opts.Start, key) <= 0 && strings.Compare(key, opts.End) < 0 {
						ch <- &common.DocInfo{Key: key, Value: val}
					}
				}
			default:
				{
					ch <- &common.DocInfo{Key: key, Value: val}
				}
			}
			return nil
		})
		close(ch)
	}()

	return ch, nil
}

// Close closes the storage
func (store *Storage) Close() error {
	return nil
}
