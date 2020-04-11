package locate

import (
	"github.com/joexu01/obj-storage/lib/rabbitmq"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate() {
	q := rabbitmq.New(os.Getenv(`RABBITMQ_SERVER`))
	defer q.Close()

	q.Bind("dataServers")
	c := q.Consume()
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
