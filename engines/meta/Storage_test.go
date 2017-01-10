package meta

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trusch/storage/testsuite"
)

type StorageSuite struct {
	testsuite.Suite
}

func TestBoltDBStorage(t *testing.T) {
	store, err := NewStorage("boltdb://test-store.db")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
	defer os.RemoveAll("./test-store.db")
}

func TestLevelDBStorage(t *testing.T) {
	store, err := NewStorage("leveldb://test-store.db")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
	defer os.RemoveAll("./test-store.db")
}

func TestMongoDBStorage(t *testing.T) {
	store, err := NewStorage("mongodb://localhost/test")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
	exec.Command("mongo", "test", "--eval", "db.dropDatabase()")
}

func TestMalformedURI(t *testing.T) {
	_, err := NewStorage("???")
	assert.Error(t, err)
	_, err = NewStorage(":")
	assert.Error(t, err)
}
