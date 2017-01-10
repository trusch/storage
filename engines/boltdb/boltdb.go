package boltdb

import (
	"bytes"

	"github.com/boltdb/bolt"
	"github.com/trusch/storage/common"
)

//Storage is an implementation for storage.Storage
type Storage struct {
	db *bolt.DB
}

// NewStorage creates a new storage instance
func NewStorage(path string) (*Storage, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, common.Error(common.InitFailed, err)
	}
	store := &Storage{db}
	return store, nil
}

// Put saves a byteslice to the db.
// Example: Save("/foo/bar", []byte{1,2,3})
func (store *Storage) Put(bucketID, key string, value []byte) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return common.Error(common.BucketNotFound)
		}
		return bucket.Put([]byte(key), value)
	})
}

// Get loads data from a path
func (store *Storage) Get(bucketID, key string) ([]byte, error) {
	var result []byte
	err := store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return common.Error(common.BucketNotFound)
		}
		key := []byte(key)
		value := bucket.Get(key)
		if value == nil {
			return common.Error(common.ReadFailed)
		}
		result = make([]byte, len(value))
		copy(result, value)
		return nil
	})
	return result, err
}

// Delete deletes a value from the db
func (store *Storage) Delete(bucketID, key string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			return common.Error(common.BucketNotFound)
		}
		return bucket.Delete([]byte(key))
	})
}

// CreateBucket creates a bucket
func (store *Storage) CreateBucket(bucketID string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketID))
		if err != nil {
			return common.Error(common.WriteFailed, err)
		}
		return nil
	})
}

// DeleteBucket deletes a bucket
func (store *Storage) DeleteBucket(bucketID string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucketID))
		if err != nil {
			return common.Error(common.WriteFailed, err)
		}
		return nil
	})
}

// List returns all Entries of a bucket
// optionally provide arguments to specifiy a key offset and a key limit
// Example: List("/foo", "abc", "xyz") -> DocInfo{Key: abc} ... DocInfo{Key: ggg} ... DocInfo{Key: xyz}
func (store *Storage) List(bucketID string, opts *common.ListOpts) (chan *common.DocInfo, error) {
	res := make(chan *common.DocInfo, 64)
	if opts == nil {
		opts = &common.ListOpts{}
	}
	earlyError := make(chan error, 2)
	go store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketID))
		if bucket == nil {
			err := common.Error(common.BucketNotFound)
			earlyError <- err
			return err
		}
		earlyError <- nil
		c := bucket.Cursor()
		switch {
		case opts.Prefix != "":
			{
				for k, v := c.Seek([]byte(opts.Prefix)); k != nil && bytes.HasPrefix(k, []byte(opts.Prefix)); k, v = c.Next() {
					res <- &common.DocInfo{Key: string(k), Value: v}
				}
			}
		case opts.Start != "":
			{
				for k, v := c.Seek([]byte(opts.Start)); k != nil && bytes.Compare(k, []byte(opts.End)) < 0; k, v = c.Next() {
					res <- &common.DocInfo{Key: string(k), Value: v}
				}
			}
		default:
			{
				for k, v := c.First(); k != nil; k, v = c.Next() {
					res <- &common.DocInfo{Key: string(k), Value: v}
				}
			}
		}
		close(res)
		return nil
	})
	err := <-earlyError
	return res, err
}

// Close closes the db
func (store *Storage) Close() error {
	err := store.db.Close()
	if err != nil {
		return common.Error(common.CloseFailed, err)
	}
	return nil
}
