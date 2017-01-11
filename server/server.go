package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/trusch/storage"
	"github.com/trusch/storage/common"
)

// Server represents the storaged webserver
type Server struct {
	store  storage.Storage
	ln     net.Listener
	server *http.Server
}

// New creates a new webserver
func New(addr string, store storage.Storage) *Server {
	srv := &http.Server{
		Addr:           addr,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server := &Server{store, nil, srv}
	server.constructRouter()
	return server
}

// ListenAndServe starts the webserver
func (srv *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", srv.server.Addr)
	if err != nil {
		return err
	}
	srv.ln = ln
	return srv.server.Serve(ln)
}

// Stop stops the webserver
func (srv *Server) Stop() error {
	err := srv.ln.Close()
	if err != nil {
		return err
	}
	return srv.store.Close()
}

func (srv *Server) constructRouter() {
	router := mux.NewRouter()
	// main ops
	router.PathPrefix("/v1/{project}/{bucket}/{key}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handlePut(w, r)
	})
	router.PathPrefix("/v1/{project}/{bucket}/{key}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleGet(w, r)
	})
	router.PathPrefix("/v1/{project}/{bucket}/{key}").Methods("DELETE").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleDelete(w, r)
	})
	router.PathPrefix("/v1/{project}/{bucket}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleCreateBucket(w, r)
	})
	router.PathPrefix("/v1/{project}/{bucket}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleList(w, r)
	})
	router.PathPrefix("/v1/{project}/{bucket}").Methods("DELETE").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleDeleteBucket(w, r)
	})
	srv.server.Handler = router
}

func (srv *Server) handlePut(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("failed put: ", r.URL.Path, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	bucket := vars["project"] + ":" + vars["bucket"]
	key := vars["key"]
	err = srv.store.Put(bucket, key, bs)
	if err != nil {
		log.Print("failed put: ", r.URL.Path, " ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (srv *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["project"] + ":" + vars["bucket"]
	key := vars["key"]
	bs, err := srv.store.Get(bucket, key)
	if err != nil {
		log.Print("failed get: ", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(bs)
}

func (srv *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["project"] + ":" + vars["bucket"]
	key := vars["key"]
	err := srv.store.Delete(bucket, key)
	if err != nil {
		log.Print("failed delete: ", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (srv *Server) handleCreateBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["project"] + ":" + vars["bucket"]
	err := srv.store.CreateBucket(bucket)
	if err != nil {
		log.Print("failed create bucket: ", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (srv *Server) handleDeleteBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["project"] + ":" + vars["bucket"]
	err := srv.store.DeleteBucket(bucket)
	if err != nil {
		log.Print("failed create bucket: ", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (srv *Server) handleList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["project"] + ":" + vars["bucket"]
	everyStr := r.FormValue("every")
	start := r.FormValue("start")
	end := r.FormValue("end")
	prefix := r.FormValue("prefix")
	var every int64
	if everyStr != "" {
		dp, e := strconv.ParseInt(everyStr, 10, 64)
		if e != nil {
			log.Print("malformed every option")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		every = dp
	}
	ch, err := srv.store.List(bucket, &common.ListOpts{Start: start, End: end, Prefix: prefix})
	if err != nil || ch == nil {
		log.Print("fail listing bucket")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if every > 0 {
		ch = reduceStream(ch, every)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("["))
	first := true
	for pair := range ch {
		if bs, err := json.Marshal(pair); err == nil {
			if first {
				first = false
			} else {
				w.Write([]byte{','})
			}
			w.Write(bs)
		}
	}
	w.Write([]byte("]"))
}

func reduceStream(input chan *common.DocInfo, every int64) chan *common.DocInfo {
	output := make(chan *common.DocInfo, 64)
	counter := int64(0)
	go func() {
		for pair := range input {
			counter++
			if counter%every == 0 {
				counter = 0
				output <- pair
			}
		}
		close(output)
	}()
	return output
}
