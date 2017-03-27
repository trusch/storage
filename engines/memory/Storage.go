package memory

import (
	"sort"
	"strings"

	"github.com/trusch/storage/common"
)

// Storage creates the apropriate store from an URI
type Storage struct {
	buckets map[string]map[string][]byte
}

// NewStorage creates a new storage from a URI
func NewStorage() (*Storage, error) {
	return &Storage{make(map[string]map[string][]byte)}, nil
}

// Put saves a byteslice to the db.
// Example: Save("/foo/bar", []byte{1,2,3})
func (store *Storage) Put(bucket, key string, value []byte) error {
	b, ok := store.buckets[bucket]
	if !ok {
		return common.Error(common.BucketNotFound)
	}
	b[key] = value
	return nil
}

// Get loads data from a key
func (store *Storage) Get(bucket, key string) ([]byte, error) {
	b, ok := store.buckets[bucket]
	if !ok {
		return nil, common.Error(common.BucketNotFound)
	}
	v, ok := b[key]
	if !ok {
		return nil, common.Error(common.ReadFailed)
	}
	return v, nil
}

// Delete deletes a value from the db
func (store *Storage) Delete(bucket, key string) error {
	b, ok := store.buckets[bucket]
	if !ok {
		return common.Error(common.BucketNotFound)
	}
	_, ok = b[key]
	if !ok {
		return nil
	}
	delete(b, key)
	return nil
}

// CreateBucket creates a bucket
func (store *Storage) CreateBucket(bucket string) error {
	if _, ok := store.buckets[bucket]; ok {
		return nil
	}
	store.buckets[bucket] = make(map[string][]byte)
	return nil
}

// DeleteBucket deletes a bucket
func (store *Storage) DeleteBucket(bucket string) error {
	if _, ok := store.buckets[bucket]; !ok {
		return common.Error(common.BucketNotFound)
	}
	delete(store.buckets, bucket)
	return nil
}

// List returns all Entries of a directory
// optionally provide arguments to specifiy a key offset and a key limit
// Example: List("/foo", "abc", "xyz") -> DocInfo{Key: abc} ... DocInfo{Key: ggg} ... DocInfo{Key: xyz}
func (store *Storage) List(bucket string, opts *common.ListOpts) (chan *common.DocInfo, error) {
	b, ok := store.buckets[bucket]
	if !ok {
		return nil, common.Error(common.BucketNotFound)
	}
	if opts == nil {
		opts = &common.ListOpts{}
	}

	keys := make([]string, len(b))
	i := 0
	for key := range b {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	ch := make(chan *common.DocInfo, 64)
	go func() {
		for _, key := range keys {
			val := b[key]
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
		}
		close(ch)
	}()

	return ch, nil
}

// Close closes the storage
func (store *Storage) Close() error {
	return nil
}
