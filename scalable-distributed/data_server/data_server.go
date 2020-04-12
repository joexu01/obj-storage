package main

import (
	"github.com/joexu01/obj-storage/scalable-distributed/data_server/heartbeat"
	"github.com/joexu01/obj-storage/scalable-distributed/data_server/locate"
	"github.com/joexu01/obj-storage/scalable-distributed/data_server/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(
		os.Getenv("LISTEN_ADDRESS"), nil))
}
