package meta

import (
	"errors"
	"net/url"
	"strings"

	"github.com/trusch/storage"
	"github.com/trusch/storage/common"
	"github.com/trusch/storage/engines/boltdb"
	"github.com/trusch/storage/engines/cache"
	"github.com/trusch/storage/engines/file"
	"github.com/trusch/storage/engines/leveldb"
	"github.com/trusch/storage/engines/memory"
	"github.com/trusch/storage/engines/mongodb"
	"github.com/trusch/storage/engines/storaged"
)

// Storage creates the apropriate store from an URI
type Storage struct {
	base storage.Storage
}

// NewStorage creates a new storage from a URI
func NewStorage(uriStr string, options ...interface{}) (*Storage, error) {
	uri, err := url.Parse(uriStr)
	if err != nil {
		return nil, err
	}
	var base storage.Storage
	switch uri.Scheme {
	case "memory":
		base, err = memory.NewStorage()
	case "cache":
		parts := strings.Split(uriStr[8:], ",")
		first, e := NewStorage(parts[0])
		if e != nil {
			return nil, e
		}
		second, e := NewStorage(parts[1])
		if e != nil {
			return nil, e
		}
		base, err = cache.NewStorage(first, second)
	case "leveldb":
		base, err = leveldb.NewStorage(uri.Host + uri.Path)
	case "boltdb":
		base, err = boltdb.NewStorage(uri.Host + uri.Path)
	case "mongodb":
		base, err = mongodb.NewStorage(uriStr)
	case "file":
		base, err = file.NewStorage(uri.Host + uri.Path)
	case "storaged":
		base, err = storaged.NewStorage(uriStr)
	case "sstoraged":
		if len(options) > 0 {
			if token, ok := options[0].(string); ok {
				base, err = storaged.NewStorage(uriStr, token)
				break
			}
		}
		base, err = storaged.NewStorage(uriStr)
	default:
		err = errors.New("unknown uri scheme, try bolt:// or leveldb://")
	}
	if err != nil {
		return nil, err
	}
	return &Storage{base}, nil
}

// Put saves a byteslice to the db.
// Example: Save("/foo/bar", []byte{1,2,3})
func (store *Storage) Put(bucket, key string, value []byte) error {
	return store.base.Put(bucket, key, value)
}

// Get loads data from a key
func (store *Storage) Get(bucket, key string) ([]byte, error) {
	return store.base.Get(bucket, key)
}

// Delete deletes a value from the db
func (store *Storage) Delete(bucket, key string) error {
	return store.base.Delete(bucket, key)
}

// CreateBucket creates a bucket
func (store *Storage) CreateBucket(bucket string) error {
	return store.base.CreateBucket(bucket)
}

// DeleteBucket deletes a bucket
func (store *Storage) DeleteBucket(bucket string) error {
	return store.base.DeleteBucket(bucket)
}

// List returns all Entries of a directory
// optionally provide arguments to specifiy a key offset and a key limit
// Example: List("/foo", "abc", "xyz") -> DocInfo{Key: abc} ... DocInfo{Key: ggg} ... DocInfo{Key: xyz}
func (store *Storage) List(bucket string, opts *common.ListOpts) (chan *common.DocInfo, error) {
	return store.base.List(bucket, opts)
}

// Close closes the storage
func (store *Storage) Close() error {
	return store.base.Close()
}
