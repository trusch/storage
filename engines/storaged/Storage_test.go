package storaged

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trusch/storage/engines/leveldb"
	"github.com/trusch/storage/server"
	"github.com/trusch/storage/testsuite"
)

type StorageSuite struct {
	testsuite.Suite
}

func TestStoragedStorage(t *testing.T) {
	baseStore, err := leveldb.NewStorage("./test-store.db")
	assert.NoError(t, err)
	server := server.New(":8080", baseStore)
	go server.ListenAndServe()
	defer server.Stop()
	store, err := NewStorage("storaged://localhost:8080/project1", "sample-token")
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
