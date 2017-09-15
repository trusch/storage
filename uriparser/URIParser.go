package uriparser

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/trusch/storage"

	"github.com/trusch/storage/base/file"
	"github.com/trusch/storage/base/gcs"
	"github.com/trusch/storage/base/mongo"

	"github.com/trusch/storage/filter/compression/gzip"
	"github.com/trusch/storage/filter/compression/snappy"
	"github.com/trusch/storage/filter/compression/xz"

	"github.com/trusch/storage/filter/encryption/aes"
	"github.com/trusch/storage/filter/encryption/ecdhe"
)

const (
	SchemeFile    = "file"
	SchemeGzip    = "gzip"
	SchemeXZ      = "xz"
	SchemeSnappy  = "snappy"
	SchemeAES     = "aes"
	SchemeMongodb = "mongodb"
	SchemeGoogle  = "google"
	SchemeECDHE   = "ecdhe"
)

// Options are arbitary key value pairs used by the various Storage impls.
type Options map[string]interface{}

// Has returns if an option exists
func (options Options) Has(key string) bool {
	_, ok := options[key]
	return ok
}

// NewFromURI creates a pipeline of Storages
// example: snappy+aes+file:///srv/backups
// -> compress, encrypt and save as files in /srv/backups
func NewFromURI(uri string, options Options) (storage.Storage, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	scheme := u.Scheme
	path := filepath.Join(u.Host, u.Path)
	schemes := strings.Split(scheme, "+")
	if options == nil {
		options = make(Options)
	}
	options["path"] = path
	options["url"] = u
	var store storage.Storage
	for i := len(schemes) - 1; i >= 0; i-- {
		s, err := createStorageByScheme(schemes[i], store, options)
		if err != nil {
			return nil, err
		}
		store = s
	}
	return store, nil
}

func createStorageByScheme(scheme string, nextStore storage.Storage, options Options) (storage.Storage, error) {
	switch scheme {
	case SchemeFile:
		{
			if path, ok := options["path"].(string); ok {
				return file.NewStorage(path), nil
			}
			return nil, errors.New("no path supplied for FileStorage")
		}
	case SchemeGoogle:
		{
			if options.Has("gcs.project") && options.Has("gcs.bucket") {
				return gcs.NewStorage(options["gcs.project"].(string), options["gcs.bucket"].(string))
			}
			if path, ok := options["path"].(string); ok {
				parts := strings.Split(path, "/")
				if len(parts) != 2 {
					return nil, errors.New("google bucket storage needs uri like google://<project-id>/<bucket-id>")
				}
				return gcs.NewStorage(parts[0], parts[1])
			}
			return nil, errors.New("no path supplied for GoogleBucketStorage")
		}
	case SchemeGzip:
		{
			level := -1
			if l, ok := options["storage.gzip.level"].(int); ok {
				level = l
			}
			return gzip.NewCompressor(nextStore, level), nil
		}
	case SchemeXZ:
		{
			return xz.NewCompressor(nextStore), nil
		}
	case SchemeSnappy:
		{
			return snappy.NewCompressor(nextStore), nil
		}
	case SchemeAES:
		{
			if key, ok := options["key"].(string); ok {
				return aes.NewEncrypter(nextStore, key), nil
			}
			return nil, errors.New("no key supplied")
		}
	case SchemeMongodb:
		{
			uri := options["url"].(*url.URL)
			uri.Scheme = "mongodb"
			db := ""
			if d, ok := options["db"].(string); ok {
				db = d
			} else if len(uri.Path) > 1 {
				db = uri.Path[1:]
			} else {
				return nil, errors.New("provide db name (either per uri, eg. mongodb://host/db or as storage option)")
			}
			return mongo.NewStorage(uri.String(), db)
		}
	case SchemeECDHE:
		{
			var (
				pubKey  string
				privKey string
				ok      bool
			)
			if pubKey, ok = options["pubkey"].(string); !ok {
				log.Warn("no pubkey supplied")
			}
			if privKey, ok = options["privkey"].(string); !ok {
				log.Warn("no privkey supplied")
			}
			if pubKey == "" && privKey == "" {
				return nil, errors.New("no key supplied")
			}
			return ecdhe.NewEncrypter(nextStore, []byte(pubKey), []byte(privKey))
		}
	}
	return nil, errors.New("unknown scheme type " + scheme)
}
