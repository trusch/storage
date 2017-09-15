package mongo

import (
	"errors"
	"io"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Storage implements the Storage interface with local files
type Storage struct {
	session *mgo.Session
	gridfs  *mgo.GridFS
}

// NewStorage creates a new filestorage with the given base directory
// if baseDirectory doesn't exist, it is created like "mkdir -p $baseDirectory"
func NewStorage(url, db string) (*Storage, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	gridfs := session.DB(db).GridFS("backup")

	return &Storage{session, gridfs}, nil
}

// GetReader returns a reader
func (fs *Storage) GetReader(id string) (io.ReadCloser, error) {
	return fs.gridfs.Open(id)
}

// GetWriter returns a writer
func (fs *Storage) GetWriter(id string) (io.WriteCloser, error) {
	return fs.gridfs.Create(id)
}

// Has returns whether an entry exists
func (fs *Storage) Has(id string) bool {
	_, err := fs.gridfs.Open(id)
	return err == nil
}

// Delete deletes an entry
func (fs *Storage) Delete(id string) error {
	if !fs.Has(id) {
		return errors.New("no such object")
	}
	return fs.gridfs.Remove(id)
}

// List returns a list of all stored objects, limited by a prefix
func (fs *Storage) List(prefix string) ([]string, error) {
	it := fs.gridfs.Files.Find(bson.M{"filename": bson.M{"$regex": "^" + prefix}}).Iter()
	res := make([]string, 0)
	for {
		row := make(map[string]interface{})
		if ok := it.Next(&row); !ok {
			break
		}
		res = append(res, row["filename"].(string))
	}
	return res, nil
}
