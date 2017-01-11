package main

import (
	"flag"
	"log"

	"github.com/mholt/caddy/server"
	"github.com/trusch/storage/engines/meta"
)

var listen = flag.String("listen", ":80", "listen address")
var backend = flag.String("backend", "leveldb:///usr/share/storaged", "backend uri")

func main() {
	flag.Parse()
	store, err := meta.NewStorage(*backend)
	if err != nil {
		log.Fatal(err)
	}
	server := server.New(*listen, store)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.ListenAndServe())
}
