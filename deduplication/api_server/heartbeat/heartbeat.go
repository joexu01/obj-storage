package heartbeat

import (
	"github.com/joexu01/obj-storage/lib/rabbitmq"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)  // 用来缓存所有数据服务节点
var mutex sync.Mutex

// ListenHeartbeat创建消息队列绑定apiServers
// exchange，聆听来自每一个节点的心跳消息
//
func ListenHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()

	q.Bind("apiServers")
	c:= q.Consume()

	go removeExpiredDataServer()
	// removeExpiredDataServer每间隔5s扫描一次dataServers
	// 清除10s内无心跳消息的数据节点

	for msg := range c {
		dataServer, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			log.Println("error parsing data server address", err)
			continue
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// GetDataServers返回当前所有数据服务节点
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

// ChooseRandomDataServer在当前所有数据服
// 务节点中选出一个
func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}
