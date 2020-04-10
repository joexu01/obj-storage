package main

import (
	"github.com/joexu01/obj-storage/introduction/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc(`/objects/`, objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
