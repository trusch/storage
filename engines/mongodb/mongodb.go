package mongodb

import (
	"github.com/trusch/storage/common"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Storage is a MongoDB implementation of the Storage interface
type Storage struct {
	session *mgo.Session
	db      *mgo.Database
}

type dbEntry struct {
	Key   string
	Value []byte
}

// NewStorage creates a new mongodb storage
func NewStorage(url string) (*Storage, error) {
	info, err := mgo.ParseURL(url)
	if err != nil {
		return nil, err
	}
	s, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	return &Storage{s, s.DB(info.Database)}, nil
}

// Put saves a byteslice to the db.
// Example: Save("/foo/bar", []byte{1,2,3})
func (store *Storage) Put(bucket, key string, value []byte) error {
	if err := store.checkBucket(bucket); err != nil {
		return err
	}
	if _, err := store.db.C(bucket).Upsert(bson.M{"key": key}, &dbEntry{key, value}); err != nil {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// Get loads data from a key
func (store *Storage) Get(bucket, key string) ([]byte, error) {
	if err := store.checkBucket(bucket); err != nil {
		return nil, err
	}
	var res dbEntry
	if err := store.db.C(bucket).Find(bson.M{"key": key}).One(&res); err != nil {
		return nil, common.Error(common.ReadFailed, err)
	}
	return res.Value, nil
}

// Delete deletes a value from the db
func (store *Storage) Delete(bucket, key string) error {
	if err := store.checkBucket(bucket); err != nil {
		return err
	}
	if err := store.db.C(bucket).Remove(bson.M{"key": key}); err != nil {
		if err != mgo.ErrNotFound {
			return common.Error(common.WriteFailed, err)
		}
	}
	return nil
}

// CreateBucket creates a bucket
func (store *Storage) CreateBucket(bucket string) error {
	if err := store.checkBucket(bucket); err == nil {
		return nil
	}
	if err := store.db.C(bucket).Insert(bson.M{}); err != nil {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// DeleteBucket deletes a bucket
func (store *Storage) DeleteBucket(bucket string) error {
	if err := store.db.C(bucket).DropCollection(); err != nil {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// List returns all Entries of a directory
// optionally provide arguments to specifiy a key offset and a key limit
// Example: List("/foo", "abc", "xyz") -> DocInfo{Key: abc} ... DocInfo{Key: ggg} ... DocInfo{Key: xyz}
func (store *Storage) List(bucket string, opts *common.ListOpts) (chan *common.DocInfo, error) {
	if err := store.checkBucket(bucket); err != nil {
		return nil, err
	}
	if opts == nil {
		opts = &common.ListOpts{}
	}
	var iter *mgo.Iter
	switch {
	case opts.Prefix != "":
		{
			iter = store.db.C(bucket).Find(bson.M{"key": bson.M{"$regex": "^" + opts.Prefix}}).Iter()
		}
	case opts.Start != "":
		{
			iter = store.db.C(bucket).Find(bson.M{"key": bson.M{"$gte": opts.Start, "$lt": opts.End}}).Iter()
		}
	default:
		{
			iter = store.db.C(bucket).Find(bson.M{"key": bson.M{"$exists": true}}).Iter()
		}
	}
	res := make(chan *common.DocInfo, 64)
	var entry dbEntry
	go func() {
		for iter.Next(&entry) {
			key := entry.Key
			val := entry.Value
			doc := &common.DocInfo{Key: key, Value: val}
			res <- doc
		}
		iter.Close()
		close(res)
	}()
	return res, nil
}

// Close closes the storage
func (store *Storage) Close() error {
	store.session.Close()
	return nil
}

func (store *Storage) checkBucket(bucket string) error {
	names, err := store.db.CollectionNames()
	if err != nil {
		return common.Error(common.ReadFailed, err)
	}
	for _, name := range names {
		if name == bucket {
			return nil
		}
	}
	return common.Error(common.BucketNotFound)
}
