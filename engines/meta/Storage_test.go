package meta

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trusch/storage/server"
	"github.com/trusch/storage/testsuite"
)

type StorageSuite struct {
	testsuite.Suite
}

func TestBoltDBStorage(t *testing.T) {
	defer os.RemoveAll("./test-store.db")
	store, err := NewStorage("boltdb://test-store.db")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}

func TestLevelDBStorage(t *testing.T) {
	defer os.RemoveAll("./test-store.db")
	store, err := NewStorage("leveldb://test-store.db")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}

func TestMongoDBStorage(t *testing.T) {
	defer exec.Command("mongo", "test", "--eval", "db.dropDatabase()")
	store, err := NewStorage("mongodb://localhost/test")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}

func TestFileStorage(t *testing.T) {
	defer os.RemoveAll("./test-store.db")
	store, err := NewStorage("file://test-store.db")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}

func TestMemoryStorage(t *testing.T) {
	store, err := NewStorage("memory://")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}

func TestStoragedStorage(t *testing.T) {
	baseStore, err := NewStorage("leveldb://test-store.db")
	assert.NoError(t, err)
	server := server.New(":8080", baseStore)
	go server.ListenAndServe()
	defer server.Stop()
	store, err := NewStorage("storaged://localhost:8080/project1")
	assert.NoError(t, err)
	defer os.RemoveAll("./test-store.db")
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}

func TestCacheStorage(t *testing.T) {
	defer os.RemoveAll("./test-store.db")
	store, err := NewStorage("cache://memory://,leveldb://test-store.db")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}

func TestMalformedURI(t *testing.T) {
	_, err := NewStorage("???")
	assert.Error(t, err)
	_, err = NewStorage(":")
	assert.Error(t, err)
}
