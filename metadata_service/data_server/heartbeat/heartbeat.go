package heartbeat

import (
	"github.com/joexu01/obj-storage/lib/rabbitmq"
	"os"
	"time"
)

// StartHeartbeat函数在一个无限循环中
// 向apiServers的exchange发送本节点的
// 监听地址
func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
