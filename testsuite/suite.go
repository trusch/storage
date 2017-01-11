package testsuite

import (
	"fmt"

	"github.com/stretchr/testify/suite"
	"github.com/trusch/storage"
	"github.com/trusch/storage/common"
)

type Suite struct {
	suite.Suite
	Store storage.Storage
}

func (suite *Suite) TestCreateUseDeleteBucket() {
	// using a unknown bucket should fail
	_, err := suite.Store.Get("bucket-name", "foo")
	suite.Error(err)
	// unknown bucket can not be deleted
	err = suite.Store.DeleteBucket("bucket-name")
	suite.Error(err)
	// creating should be possible
	err = suite.Store.CreateBucket("bucket-name")
	suite.NoError(err)
	// created bucket should be writeable
	err = suite.Store.Put("bucket-name", "foo", []byte("hello world"))
	suite.NoError(err)
	// created bucket should be readable
	val, err := suite.Store.Get("bucket-name", "foo")
	suite.NoError(err)
	suite.Equal("hello world", string(val))
	// created bucket can be deleted
	err = suite.Store.DeleteBucket("bucket-name")
	suite.NoError(err)
	// using a deleted bucket should fail
	_, err = suite.Store.Get("bucket-name", "foo")
	suite.Error(err)
}

func (suite *Suite) TestListAll() {
	err := suite.Store.CreateBucket("bucket-name")
	suite.NoError(err)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%03d", i) // 001, 002, ... 010 ... 100
		val := []byte(key)
		err = suite.Store.Put("bucket-name", key, val)
		suite.NoError(err)
	}
	ch, err := suite.Store.List("bucket-name", nil)
	suite.NoError(err)
	for i := 0; i < 100; i++ {
		info, ok := <-ch
		suite.True(ok)
		if !ok {
			suite.FailNow("no data")
		}
		expectedKey := fmt.Sprintf("%03d", i)
		expectedVal := expectedKey
		suite.Equal(expectedKey, info.Key)
		suite.Equal(expectedVal, string(info.Value))
	}
	_, ok := <-ch
	suite.False(ok)
	err = suite.Store.DeleteBucket("bucket-name")
	suite.NoError(err)
}

func (suite *Suite) TestListPrefix() {
	err := suite.Store.CreateBucket("bucket-name")
	suite.NoError(err)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%03d", i) // 001, 002, ... 010 ... 100
		val := []byte(key)
		err = suite.Store.Put("bucket-name", key, val)
		suite.NoError(err)
	}
	ch, err := suite.Store.List("bucket-name", &common.ListOpts{Prefix: "01"})
	suite.NoError(err)
	for i := 10; i < 20; i++ {
		info, ok := <-ch
		suite.True(ok)
		if !ok {
			suite.FailNow("no data")
		}
		expectedKey := fmt.Sprintf("%03d", i)
		expectedVal := expectedKey
		suite.Equal(expectedKey, info.Key)
		suite.Equal(expectedVal, string(info.Value))
	}
	_, ok := <-ch
	suite.False(ok)
	err = suite.Store.DeleteBucket("bucket-name")
	suite.NoError(err)
}

func (suite *Suite) TestListRange() {
	err := suite.Store.CreateBucket("bucket-name")
	suite.NoError(err)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%03d", i) // 001, 002, ... 010 ... 099
		val := []byte(key)
		err = suite.Store.Put("bucket-name", key, val)
		suite.NoError(err)
	}
	ch, err := suite.Store.List("bucket-name", &common.ListOpts{Start: "023", End: "100"})
	suite.NoError(err)
	for i := 23; i < 100; i++ {
		info, ok := <-ch
		suite.True(ok)
		if !ok {
			suite.FailNow("no data")
		}
		expectedKey := fmt.Sprintf("%03d", i)
		expectedVal := expectedKey
		suite.Equal(expectedKey, info.Key)
		suite.Equal(expectedVal, string(info.Value))
	}
	_, ok := <-ch
	suite.False(ok)
	err = suite.Store.DeleteBucket("bucket-name")
	suite.NoError(err)
}

func (suite *Suite) TestPutInNonExistingBucket() {
	err := suite.Store.Put("unknown", "foo", []byte("hello"))
	suite.Error(err)
}

func (suite *Suite) TestListNonExistingBucket() {
	_, err := suite.Store.List("bucket-name", nil)
	suite.Error(err)
}

func (suite *Suite) TestGetNonExistingKey() {
	err := suite.Store.CreateBucket("bucket-name")
	suite.NoError(err)
	_, err = suite.Store.Get("bucket-name", "foo")
	suite.Error(err)
	err = suite.Store.DeleteBucket("bucket-name")
	suite.NoError(err)
}

func (suite *Suite) TestDelete() {
	err := suite.Store.CreateBucket("bucket-name")
	suite.NoError(err)
	err = suite.Store.Put("bucket-name", "foo", []byte("hello"))
	suite.NoError(err)
	err = suite.Store.Delete("bucket-name", "foo")
	suite.NoError(err)
	err = suite.Store.Delete("bucket-name", "foo")
	suite.NoError(err) // (!)
	err = suite.Store.Delete("wrong", "foo")
	suite.Error(err)
	err = suite.Store.DeleteBucket("bucket-name")
	suite.NoError(err)
}

func (suite *Suite) TestRecreateBucket() {
	err := suite.Store.CreateBucket("bucket-name")
	suite.NoError(err)
	err = suite.Store.CreateBucket("bucket-name")
	suite.NoError(err)
	err = suite.Store.DeleteBucket("bucket-name")
	suite.NoError(err)
}
