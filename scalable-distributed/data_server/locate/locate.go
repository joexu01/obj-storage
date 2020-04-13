package locate

import (
	"github.com/joexu01/obj-storage/lib/rabbitmq"
	"os"
	"strconv"
)

// Locate函数访问磁盘上对应的文件名
// 并判断文件是否存在
func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// StartLocate函数创建RabbitMQ结构体
// 绑定dataServers的exchange
// 从go channel接收的消息的正文是需要
// 定位的对象的名称，若存在则返回监听地址
func StartLocate() {
	q := rabbitmq.New(os.Getenv(`RABBITMQ_SERVER`))
	defer q.Close()

	q.Bind("dataServers")
	c := q.Consume()
	// 遍历返回的go channel
	for msg := range c {
		obj, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + obj) {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
