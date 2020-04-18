package objects

import (
	"io"
	"log"
	"net/http"
)

func storeObject(r io.Reader, object string) (respCode int, err error) {
	stream, err := putStream(object)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	_, err = io.Copy(stream, r)
	if err != nil {
		log.Println("object storage error: ", err)
	}
	err = stream.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
