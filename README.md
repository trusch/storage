storage
=======

Storage is a library and daemon that abstracts different storage engines away.
It provides basic Key/Value Storage functionality and works from embedded devices up to big cloud services.

### Core Methods

* Create bucket
* Delete bucket
* Put value into bucket
* Get value from bucket
* Delete value from bucket
* Get all values from bucket
* Get all values with specific key prefix from bucket
* Get all values within a specific key range from bucket

### Supported Engines

* LevelDB
* MongoDB
* BoltDB
* File-based
* Memory
* Storaged
* Cache (combine two other storage engines)

### API Server

github.com/trusch/storage/storaged contains a daemon which provides the core methods via HTTP.

#### Installation
Install it via `go get github.com/trusch/storage/storaged`

#### Bucket Management

* Create Bucket
  * `PUT /v1/my-project/my-bucket`
* Delete Bucket
  * `DELETE /v1/my-project/my-bucket`

#### Data Management

* Put value
  * `PUT /v1/my-project/my-bucket/my-key`
  * Complete HTTP body is threated as value
* Get value
  * `GET /v1/my-project/my-bucket/my-key`
* Delete value
  * `DELETE /v1/my-project/my-bucket/my-key`
* List all values
  * `GET /v1/my-project/my-bucket`
* List values with key prefix
  * `GET /v1/my-project/my-bucket?prefix=abc`
* List values within range
  * `GET /v1/my-project/my-bucket?start=abc&end=xyz`
  * start is inclusive, end not

### Code Example
```go
package main

import (
  "log"

  "github.com/trusch/storage/meta"
)

func main(){
  store, err := meta.NewStorage("leveldb://./test-store.db")
  if err != nil {
    log.Fatal(err)
  }
  err = store.CreateBucket("my-bucket")
  if err != nil {
    log.Fatal(err)
  }
  err = store.Put("my-bucket", "my-key", []byte("hello world"))
  if err != nil {
    log.Fatal(err)
  }
  val, err := store.Get("my-bucket", "my-key")
  if err != nil {
    log.Fatal(err)
  }
  if val != []byte("hello world") {
    log.Fatal("wrong result, this will not happen ;)")
  }
}

```
