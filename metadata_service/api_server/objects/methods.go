package objects

import (
	"github.com/joexu01/obj-storage/lib/es"
	"github.com/joexu01/obj-storage/lib/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), `/`)[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var err error

	if len(versionId) != 0 {
		version, err = strconv.Atoi(versionId[0])
		if err != nil {
			log.Println("error handling get request:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	meta, err := es.GetMetadata(name, version)
	if err != nil {
		log.Println("error handling get request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	object := url.PathEscape(meta.Hash)
	stream, err := getStream(object)
	if err != nil {
		log.Println("error handling get request:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = io.Copy(w, stream)
	if err != nil {
		log.Println("error copying file:", err)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("error handling put method: missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respCode, err := storeObject(r.Body, url.PathEscape(hash))
	if err != nil {
		log.Println("error handling put request:", err)
		w.WriteHeader(respCode)
		return
	}
	if respCode != http.StatusOK {
		w.WriteHeader(respCode)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), `/`)[2]
	size := utils.GetSizeFromHeader(r.Header)
	err = es.AddVersion(name, hash, size)
	if err != nil {
		log.Println("error handling put request:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func del(w http.ResponseWriter, r *http.Request) {

}
