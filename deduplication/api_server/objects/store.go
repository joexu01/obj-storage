package objects

import (
	"fmt"
	"github.com/joexu01/obj-storage/deduplication/api_server/locate"
	"github.com/joexu01/obj-storage/lib/utils"
	"io"
	"log"
	"net/http"
	"net/url"
)

func storeObject(r io.Reader, hash string, size int64) (respCode int, err error) {
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, err := putStream(url.PathEscape(hash), size)
	if err != nil {
		log.Println("error storing object:", err)
		return http.StatusInternalServerError, err
	}

	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf(
			"object hash mismatch, calculated=%s but requested=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
