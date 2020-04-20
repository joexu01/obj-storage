package main

import (
	"github.com/joexu01/obj-storage/metadata_service/api_server/heartbeat"
	"github.com/joexu01/obj-storage/metadata_service/api_server/locate"
	"github.com/joexu01/obj-storage/metadata_service/api_server/objects"
	"github.com/joexu01/obj-storage/metadata_service/api_server/versions"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
