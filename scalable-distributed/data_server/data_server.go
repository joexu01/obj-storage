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
	go heartbeat.StartHeartbeat()  // 不断发送本节点地址
	go locate.StartLocate()  // 文件定位服务
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(
		os.Getenv("LISTEN_ADDRESS"), nil))
}
