package leveldb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trusch/storage/testsuite"
)

type StorageSuite struct {
	testsuite.Suite
}

func (suite *StorageSuite) TestFailingInit() {
	_, err := NewStorage("/root/db")
	suite.Error(err)
	_, err = NewStorage("/doesnt-exist")
	suite.Error(err)
}

func (suite *StorageSuite) TestClose() {
	store, err := NewStorage("./close-test.db")
	suite.NoError(err)
	err = store.Close()
	suite.NoError(err)
	err = store.Close()
	suite.NoError(err) // (!)
	os.RemoveAll("./close-test.db")
}

func TestLevelDBStorage(t *testing.T) {
	store, err := NewStorage("./test-store.db")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	defer os.RemoveAll("./test-store.db")
}
