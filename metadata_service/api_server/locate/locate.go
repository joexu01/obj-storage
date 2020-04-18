package locate

import (
	"github.com/joexu01/obj-storage/lib/rabbitmq"
	"os"
	"strconv"
	"time"
)

// Locate函数接受一个需要定位的对象的名称
// 向data servers的exchange群发该名称
// 1s内无响应则关闭消息队列
func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <- c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

func Exist(name string) bool {
	return Locate(name) != ""
}
