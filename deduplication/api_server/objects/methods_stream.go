package objects

import (
	"fmt"
	"github.com/joexu01/obj-storage/deduplication/api_server/heartbeat"
	"github.com/joexu01/obj-storage/deduplication/api_server/locate"
	"github.com/joexu01/obj-storage/lib/objectstream"
	"io"
)

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf(
			"failed to locate object %s", object)
	}
	return objectstream.NewGetStream(server, object)
}

func putStream(hash string, size int64) (*objectstream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("failed to find any data server")
	}
	return objectstream.NewTempPutStream(server, hash, size)
}
