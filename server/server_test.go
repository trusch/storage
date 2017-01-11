package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/trusch/storage/common"
	"github.com/trusch/storage/engines/meta"
)

type ServerSuite struct {
	suite.Suite
	srv *Server
}

func (suite *ServerSuite) SetupSuite() {
	store, err := meta.NewStorage("leveldb://test-store.db")
	suite.NoError(err)
	suite.NotEmpty(store)
	suite.srv = New(":8080", store)
	go suite.srv.ListenAndServe()
	time.Sleep(200 * time.Millisecond)
}

func (suite *ServerSuite) TearDownTest() {
	err := suite.srv.store.Close()
	suite.NoError(err)
	err = os.RemoveAll("./test-store.db")
	suite.NoError(err)
	store, err := meta.NewStorage("leveldb://test-store.db")
	suite.NoError(err)
	suite.NotEmpty(store)
	suite.srv.store = store
}

func (suite *ServerSuite) TearDownSuite() {
	err := suite.srv.Stop()
	suite.NoError(err)
	err = os.RemoveAll("./test-store.db")
	suite.NoError(err)
}

func (suite *ServerSuite) TestPut() {
	res, err := suite.request("PUT", "/p1/mybucket/foo", "hello world")
	suite.Error(err)
	suite.Empty(res)
	res, err = suite.request("PUT", "/p1/mybucket", "")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("PUT", "/p1/mybucket/foo", "hello world")
	suite.NoError(err)
	suite.Empty(res)
}

func (suite *ServerSuite) TestGet() {
	res, err := suite.request("PUT", "/p1/mybucket", "")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("PUT", "/p1/mybucket/foo", "hello world")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("GET", "/p1/mybucket/foo", "")
	suite.NoError(err)
	suite.Equal("hello world", res)
}

func (suite *ServerSuite) TestDelete() {
	res, err := suite.request("PUT", "/p1/mybucket", "")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("PUT", "/p1/mybucket/foo", "hello world")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("DELETE", "/p1/mybucket/foo", "")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("GET", "/p1/mybucket/foo", "")
	suite.Equal("404", err.Error())
	suite.Empty(res)
}

func (suite *ServerSuite) TestDeleteBucket() {
	_, err := suite.request("DELETE", "/p1/mybucket", "")
	suite.Error(err)
	_, err = suite.request("PUT", "/p1/mybucket", "")
	suite.NoError(err)
	_, err = suite.request("DELETE", "/p1/mybucket", "")
	suite.NoError(err)
	_, err = suite.request("DELETE", "/p1/mybucket", "")
	suite.Error(err)
}

func (suite *ServerSuite) TestCreateBucket() {
	_, err := suite.request("PUT", "/p1/mybucket", "")
	suite.NoError(err)
	_, err = suite.request("PUT", "/p1/mybucket", "")
	suite.NoError(err)
	_, err = suite.request("DELETE", "/p1/mybucket", "")
	suite.NoError(err)
	_, err = suite.request("DELETE", "/p1/mybucket", "")
	suite.Error(err)
}

func (suite *ServerSuite) TestGetRange() {
	res, err := suite.request("PUT", "/p1/mybucket", "")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("PUT", "/p1/mybucket/0", "0")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("PUT", "/p1/mybucket/1", "1")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("PUT", "/p1/mybucket/2", "2")
	suite.NoError(err)
	suite.Empty(res)
	res, err = suite.request("PUT", "/p1/mybucket/3", "3")
	suite.NoError(err)
	suite.Empty(res)

	res, err = suite.request("GET", "/p1/mybucket", "")
	suite.NoError(err)
	slice := make([]*common.DocInfo, 0)
	decoder := json.NewDecoder(strings.NewReader(res))
	err = decoder.Decode(&slice)
	suite.NoError(err)
	suite.Equal(4, len(slice))
	suite.Equal("0", string(slice[0].Value))
	suite.Equal("1", string(slice[1].Value))
	suite.Equal("2", string(slice[2].Value))
	suite.Equal("3", string(slice[3].Value))
}

func (suite *ServerSuite) TestGetRangeWithEvery() {
	var (
		count = 1000
		every = 100
	)
	res, err := suite.request("PUT", "/p1/mybucket", "")
	suite.NoError(err)
	suite.Empty(res)
	for i := 0; i < count; i++ {
		res, err = suite.request("PUT", fmt.Sprintf("/p1/mybucket/%v", i), fmt.Sprintf("%v", i))
		suite.NoError(err)
		suite.Empty(res)
	}
	res, err = suite.request("GET", fmt.Sprintf("/p1/mybucket?every=%v", every), "")
	suite.NoError(err)
	slice := make([]map[string]interface{}, 0)
	decoder := json.NewDecoder(strings.NewReader(res))
	err = decoder.Decode(&slice)
	suite.NoError(err)
	suite.True(len(slice) == count/every)
}

func (suite *ServerSuite) request(method, path string, data string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:8080/v1%v", path), strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return string(body), fmt.Errorf("%v", resp.StatusCode)
	}
	return string(body), nil
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
