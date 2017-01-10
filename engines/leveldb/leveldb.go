package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/trusch/storage/common"
)

// Storage is a leveldb implementation of the storage interface
type Storage struct {
	db *leveldb.DB
}

// NewStorage opens a new leveldb database
func NewStorage(path string) (*Storage, error) {
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}
	db, err := leveldb.OpenFile(path, o)
	if err != nil {
		return nil, err
	}
	store := &Storage{db}
	return store, nil
}

// Put saves a byteslice to the db.
// Example: Save("/foo/bar", []byte{1,2,3})
func (store *Storage) Put(bucket, key string, value []byte) error {
	if err := store.checkBucket(bucket); err != nil {
		return err
	}
	err := store.db.Put([]byte(bucket+"/"+key), value, nil)
	if err != nil {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// Get loads data from a key
func (store *Storage) Get(bucket, key string) ([]byte, error) {
	val, err := store.db.Get([]byte(bucket+"/"+key), nil)
	if err != nil {
		return nil, common.Error(common.ReadFailed, err)
	}
	return val, nil
}

// Delete deletes a value from the db
func (store *Storage) Delete(bucket, key string) error {
	if err := store.checkBucket(bucket); err != nil {
		return err
	}
	err := store.db.Delete([]byte(bucket+"/"+key), nil)
	if err != nil {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// CreateBucket creates a bucket
func (store *Storage) CreateBucket(bucket string) error {
	err := store.db.Put([]byte(bucket), []byte{}, nil)
	if err != nil {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// DeleteBucket deletes a bucket
func (store *Storage) DeleteBucket(bucket string) error {
	if err := store.checkBucket(bucket); err != nil {
		return err
	}
	ch, err := store.List(bucket, nil)
	if err != nil {
		return err
	}
	for item := range ch {
		store.Delete(bucket, item.Key)
	}
	err = store.db.Delete([]byte(bucket), nil)
	if err != nil {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// List returns all Entries of a directory
// optionally provide arguments to specifiy a key offset and a key limit
// Example: List("/foo", "abc", "xyz") -> DocInfo{Key: abc} ... DocInfo{Key: ggg} ... DocInfo{Key: xyz}
func (store *Storage) List(bucket string, opts *common.ListOpts) (chan *common.DocInfo, error) {
	err := store.checkBucket(bucket)
	if err != nil {
		return nil, err
	}
	if opts == nil {
		opts = &common.ListOpts{}
	}
	var iter iterator.Iterator
	switch {
	case opts.Prefix != "":
		{
			iter = store.db.NewIterator(util.BytesPrefix([]byte(bucket+"/"+opts.Prefix)), nil)
		}
	case opts.Start != "":
		{
			iter = store.db.NewIterator(&util.Range{Start: []byte(bucket + "/" + opts.Start), Limit: []byte(bucket + "/" + opts.End)}, nil)
		}
	default:
		{
			iter = store.db.NewIterator(util.BytesPrefix([]byte(bucket+"/")), nil)
		}
	}
	res := make(chan *common.DocInfo, 64)
	bucketNameLen := len(bucket) + 1
	go func() {
		for iter.Next() {
			key := string(iter.Key()[bucketNameLen:])
			val := iter.Value()
			doc := &common.DocInfo{Key: key, Value: make([]byte, len(val))}
			copy(doc.Value, val)
			res <- doc
		}
		iter.Release()
		close(res)
	}()
	return res, nil
}

// Close closes the storage
func (store *Storage) Close() error {
	err := store.db.Close()
	if err != nil {
		if err == leveldb.ErrClosed {
			return nil
		}
		return common.Error(common.CloseFailed, err)
	}
	return nil
}

func (store *Storage) checkBucket(bucket string) error {
	ok, err := store.db.Has([]byte(bucket), nil)
	if err != nil {
		return common.Error(common.ReadFailed, err)
	}
	if !ok {
		return common.Error(common.BucketNotFound)
	}
	return nil
}
