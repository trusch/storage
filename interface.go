package storage

import "github.com/trusch/storage/common"

// Storage specifies the commandset for each storage
type Storage interface {
	// Put saves a byteslice to the db.
	// Example: Save("/foo/bar", []byte{1,2,3})
	Put(bucket, key string, value []byte) error
	// Get loads data from a key
	Get(bucket, key string) ([]byte, error)
	// Delete deletes a value from the db
	Delete(bucket, key string) error
	// CreateBucket creates a bucket
	CreateBucket(bucket string) error
	// DeleteBucket deletes a bucket
	DeleteBucket(bucket string) error
	// List returns all Entries of a directory
	// optionally provide arguments to specifiy a key offset and a key limit
	// Example: List("/foo", "abc", "xyz") -> DocInfo{Key: abc} ... DocInfo{Key: ggg} ... DocInfo{Key: xyz}
	List(bucket string, opts *common.ListOpts) (chan *common.DocInfo, error)
	// Close closes the storage
	Close() error
}
