package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Storage implements the Storage interface with local files
type Storage struct {
	baseDirectory string
}

// NewStorage creates a new filestorage with the given base directory
// if baseDirectory doesn't exist, it is created like "mkdir -p $baseDirectory"
func NewStorage(baseDirectory string) *Storage {
	baseDirectory = filepath.Clean(baseDirectory)
	os.MkdirAll(baseDirectory, 0700)
	return &Storage{baseDirectory}
}

// GetReader returns a reader
func (fs *Storage) GetReader(id string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(fs.baseDirectory, id))
}

// GetWriter returns a writer
func (fs *Storage) GetWriter(id string) (io.WriteCloser, error) {
	return os.Create(filepath.Join(fs.baseDirectory, id))
}

// Has returns whether an entry exists
func (fs *Storage) Has(id string) bool {
	_, err := os.Stat(filepath.Join(fs.baseDirectory, id))
	return err == nil
}

// Delete deletes an entry
func (fs *Storage) Delete(id string) error {
	return os.Remove(filepath.Join(fs.baseDirectory, id))
}

// List returns a list of all stored objects, limited by a prefix
func (fs *Storage) List(prefix string) ([]string, error) {
	files, err := ioutil.ReadDir(fs.baseDirectory)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)
	for _, f := range files {
		name := f.Name()
		if strings.HasPrefix(name, prefix) {
			res = append(res, name)
		}
	}
	return res, nil
}
