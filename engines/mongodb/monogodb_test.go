package mongodb

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trusch/storage/testsuite"
)

type StorageSuite struct {
	testsuite.Suite
}

func (suite *StorageSuite) TestClose() {
	store, err := NewStorage("localhost")
	suite.NoError(err)
	err = store.Close()
	suite.NoError(err)
	err = store.Close()
	suite.NoError(err) // (!)
	exec.Command("mongo", "test", "--eval", "db.dropDatabase()")
}

func TestMongoDBStorage(t *testing.T) {
	store, err := NewStorage("localhost")
	assert.NoError(t, err)
	s := &StorageSuite{}
	s.Store = store
	suite.Run(t, s)
	exec.Command("mongo", "test", "--eval", "db.dropDatabase()")
}
