package objects

import (
	"fmt"
	"github.com/joexu01/obj-storage/lib/objectstream"
	"github.com/joexu01/obj-storage/scalable-distributed/api_server/locate"
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
