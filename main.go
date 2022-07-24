package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func newServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", fileHandler)

	return &http.Server{
		Addr:        addr,
		Handler:     mux,
		ReadTimeout: 2 * time.Second,
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI[1:]
	if len(path) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f, err := os.Open("./" + path)
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if os.IsPermission(err) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer := make([]byte, 4096)

	for {
		rn, err := f.Read(buffer)
		if err == io.EOF {
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cn, err := w.Write(buffer[:rn])
		if err != nil || cn != rn {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	port := flag.Int("port", 3333, "port")
	host := flag.String("host", "0.0.0.0", "host")
	flag.Parse()

	addr := net.JoinHostPort(*host, strconv.Itoa(*port))
	srv := newServer(addr)
	log.Fatal(srv.ListenAndServe())
}
