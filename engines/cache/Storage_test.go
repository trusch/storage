package cache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trusch/storage/engines/leveldb"
	"github.com/trusch/storage/engines/memory"
	"github.com/trusch/storage/testsuite"
)

type StorageSuite struct {
	testsuite.Suite
}

func TestBoltDBStorage(t *testing.T) {
	defer os.RemoveAll("./test-store.db")
	first, err := memory.NewStorage()
	assert.NoError(t, err)
	second, err := leveldb.NewStorage("./test-store.db")
	assert.NoError(t, err)
	store, err := NewStorage(first, second)
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	err = store.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}
