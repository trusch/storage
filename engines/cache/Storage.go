package cache

import (
	"errors"

	"github.com/trusch/storage"
	"github.com/trusch/storage/common"
)

// Storage creates the apropriate store from an URI
type Storage struct {
	first  storage.Storage
	second storage.Storage
}

// NewStorage creates a new storage from a URI
func NewStorage(first, second storage.Storage) (*Storage, error) {
	return &Storage{first, second}, nil
}

// Put saves a byteslice to the db.
// Example: Save("/foo/bar", []byte{1,2,3})
func (store *Storage) Put(bucket, key string, value []byte) error {
	if err := store.first.Put(bucket, key, value); err != nil {
		return common.Error(common.WriteFailed, errors.New("first level fail"), err)
	}
	if err := store.second.Put(bucket, key, value); err != nil {
		return common.Error(common.WriteFailed, errors.New("second level fail"), err)
	}
	return nil
}

// Get loads data from a key
func (store *Storage) Get(bucket, key string) ([]byte, error) {
	val, err := store.first.Get(bucket, key)
	if err != nil {
		val, err = store.second.Get(bucket, key)
		if err == nil {
			store.first.Put(bucket, key, val)
		}
	}
	return val, err
}

// Delete deletes a value from the db
func (store *Storage) Delete(bucket, key string) error {
	if err := store.first.Delete(bucket, key); err != nil {
		return common.Error(common.WriteFailed, errors.New("first level fail"), err)
	}
	if err := store.second.Delete(bucket, key); err != nil {
		return common.Error(common.WriteFailed, errors.New("second level fail"), err)
	}
	return nil
}

// CreateBucket creates a bucket
func (store *Storage) CreateBucket(bucket string) error {
	if err := store.first.CreateBucket(bucket); err != nil {
		return common.Error(common.WriteFailed, errors.New("first level fail"), err)
	}
	if err := store.second.CreateBucket(bucket); err != nil {
		return common.Error(common.WriteFailed, errors.New("second level fail"), err)
	}
	return nil
}

// DeleteBucket deletes a bucket
func (store *Storage) DeleteBucket(bucket string) error {
	if err := store.first.DeleteBucket(bucket); err != nil {
		return common.Error(common.WriteFailed, errors.New("first level fail"), err)
	}
	if err := store.second.DeleteBucket(bucket); err != nil {
		return common.Error(common.WriteFailed, errors.New("second level fail"), err)
	}
	return nil
}

// List returns all Entries of a directory
// optionally provide arguments to specifiy a key offset and a key limit
// Example: List("/foo", "abc", "xyz") -> DocInfo{Key: abc} ... DocInfo{Key: ggg} ... DocInfo{Key: xyz}
func (store *Storage) List(bucket string, opts *common.ListOpts) (chan *common.DocInfo, error) {
	ch, err := store.first.List(bucket, opts)
	if err != nil {
		ch, err = store.second.List(bucket, opts)
		if err != nil {
			return nil, err
		}
	}
	return ch, nil
}

// Close closes the storage
func (store *Storage) Close() error {
	err1 := store.first.Close()
	err2 := store.second.Close()
	if err1 != nil && err2 == nil {
		return common.Error(common.CloseFailed, errors.New("first level fail"), err1)
	}
	if err1 == nil && err2 != nil {
		return common.Error(common.CloseFailed, errors.New("second level fail"), err2)
	}
	if err1 != nil && err2 != nil {
		return common.Error(common.CloseFailed, errors.New("both levels failed"), err1, err2)
	}
	return nil
}
