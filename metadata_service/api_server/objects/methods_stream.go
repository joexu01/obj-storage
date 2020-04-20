package objects

import (
	"fmt"
	"github.com/joexu01/obj-storage/lib/objectstream"
	"github.com/joexu01/obj-storage/metadata_service/api_server/heartbeat"
	"github.com/joexu01/obj-storage/metadata_service/api_server/locate"
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

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("failed to find any data server")
	}
	return objectstream.NewPutStream(server, object), nil
}
