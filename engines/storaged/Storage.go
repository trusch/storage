package storaged

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/trusch/storage/common"
)

// Storage creates the apropriate store from an URI
type Storage struct {
	client  *http.Client
	baseURL string
	token   string
}

// NewStorage creates a new storage from a URI
func NewStorage(uriStr string, token ...string) (*Storage, error) {
	uri, err := url.Parse(uriStr)
	if err != nil {
		return nil, err
	}
	base := ""
	switch uri.Scheme {
	case "storaged":
		{
			base = fmt.Sprintf("http://%v/v1%v", uri.Host, uri.Path)
		}
	case "sstoraged":
		{
			base = fmt.Sprintf("https://%v/v1%v", uri.Host, uri.Path)
		}
	}
	t := ""
	if len(token) > 0 {
		t = token[0]
	}
	return &Storage{&http.Client{}, base, t}, nil
}

// Put saves a byteslice to the db.
// Example: Save("/foo/bar", []byte{1,2,3})
func (store *Storage) Put(bucket, key string, value []byte) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%v/%v/%v", store.baseURL, bucket, key), bytes.NewReader(value))
	if err != nil {
		return common.Error(common.WriteFailed, err)
	}
	if store.token != "" {
		req.Header.Set("Autorization", fmt.Sprintf("bearer %v", store.token))
	}
	resp, err := store.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// Get loads data from a key
func (store *Storage) Get(bucket, key string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/%v/%v", store.baseURL, bucket, key), nil)
	if err != nil {
		return nil, common.Error(common.ReadFailed, err)
	}
	if store.token != "" {
		req.Header.Set("Autorization", fmt.Sprintf("bearer %v", store.token))
	}
	resp, err := store.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, common.Error(common.ReadFailed, err)
	}
	val, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, common.Error(common.ReadFailed, err)
	}
	return val, nil
}

// Delete deletes a value from the db
func (store *Storage) Delete(bucket, key string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/%v/%v", store.baseURL, bucket, key), nil)
	if err != nil {
		return common.Error(common.WriteFailed, err)
	}
	if store.token != "" {
		req.Header.Set("Autorization", fmt.Sprintf("bearer %v", store.token))
	}
	resp, err := store.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// CreateBucket creates a bucket
func (store *Storage) CreateBucket(bucket string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%v/%v", store.baseURL, bucket), nil)
	if err != nil {
		return common.Error(common.WriteFailed, err)
	}
	if store.token != "" {
		req.Header.Set("Autorization", fmt.Sprintf("bearer %v", store.token))
	}
	resp, err := store.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// DeleteBucket deletes a bucket
func (store *Storage) DeleteBucket(bucket string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/%v", store.baseURL, bucket), nil)
	if err != nil {
		return common.Error(common.WriteFailed, err)
	}
	if store.token != "" {
		req.Header.Set("Autorization", fmt.Sprintf("bearer %v", store.token))
	}
	resp, err := store.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return common.Error(common.WriteFailed, err)
	}
	return nil
}

// List returns all Entries of a directory
// optionally provide arguments to specifiy a key offset and a key limit
// Example: List("/foo", "abc", "xyz") -> DocInfo{Key: abc} ... DocInfo{Key: ggg} ... DocInfo{Key: xyz}
func (store *Storage) List(bucket string, opts *common.ListOpts) (chan *common.DocInfo, error) {
	uri := fmt.Sprintf("%v/%v", store.baseURL, bucket)
	if opts == nil {
		opts = &common.ListOpts{}
	}
	if opts.Prefix != "" {
		uri = fmt.Sprintf("%v?prefix=%v", uri, opts.Prefix)
	} else if opts.Start != "" {
		uri = fmt.Sprintf("%v?start=%v&end=%v", uri, opts.Start, opts.End)
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, common.Error(common.ReadFailed, err)
	}
	if store.token != "" {
		req.Header.Set("Autorization", fmt.Sprintf("bearer %v", store.token))
	}
	resp, err := store.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, common.Error(common.ReadFailed, err)
	}
	ch := make(chan *common.DocInfo, 64)
	go func() {
		defer close(ch)
		reader := bufio.NewReader(resp.Body)
		b, err := reader.ReadByte()
		if err != nil || b != '[' {
			return
		}
		for {
			next, err := reader.ReadSlice('}')
			if err != nil {
				log.Print(err)
				return
			}
			info := &common.DocInfo{}
			err = json.Unmarshal(next, info)
			if err != nil {
				log.Print(err)
				return
			}
			ch <- info
			b, err := reader.ReadByte()
			if err != nil || b != ',' {
				return
			}
		}
	}()
	return ch, nil
}

// Close closes the storage
func (store *Storage) Close() error {
	return nil
}
