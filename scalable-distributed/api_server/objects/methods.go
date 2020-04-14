package objects

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := getStream(object)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = io.Copy(w, stream)
	if err != nil {
		log.Println("error occurred when handling a GET request")
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, e := storeObject(r.Body, object)
	if e != nil {
		log.Println(e)
	}
	w.WriteHeader(c)
}
