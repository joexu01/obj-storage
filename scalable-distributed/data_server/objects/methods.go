package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(os.Getenv("STORAGE_ROOT")+"/objects/"+
		strings.Split(r.URL.EscapedPath(), "/")[2])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Println(err)
		}
	}()

	_, err = io.Copy(w, file)
	if err != nil {
		log.Println(err)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	file, err := os.Create(os.Getenv("STORAGE_ROOT") +  `/objects/` +
		strings.Split(r.URL.EscapedPath(), "/")[2])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Println(closeErr)
		}
	}()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Println(err)
	}
}
