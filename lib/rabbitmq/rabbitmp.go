package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	channel  *amqp.Channel
	conn     *amqp.Connection
	Name     string
	exchange string // 交换器用于接收、分配消息
	// reference: https://www.cnblogs.com/javaGoGo/p/10111513.html
}

func New(s string) *RabbitMQ {
	// 建立连接并开启Advanced Message Queuing Protocol信道
	conn, err := amqp.Dial(s)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	queue, err := ch.QueueDeclare("", false, true,
		false, false, nil)
	if err != nil {
		panic(err)
	}

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.conn = conn
	mq.Name = queue.Name
	return mq
}

// Bind方法将接收者的消息队列和一个exchange绑定所有
// 发往该exchange的消息都会在自己的消息队列中接收到
func (q *RabbitMQ) Bind(exchange string) {
	err := q.channel.QueueBind(q.Name, "", exchange, false, nil)
	if err != nil {
		panic(err)
	}
	q.exchange = exchange
}

// Send方法可以向某个消息队列发送消息
func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	err = q.channel.Publish("", queue, false, false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    str,
		})
	if err != nil {
		panic(err)
	}
}

// Publish方法可以向某个exchange发送消息
func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	err = q.channel.Publish(exchange, "", false, false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    str,
		})
	if err != nil {
		panic(err)
	}
}

// Consume方法生成一个接收消息的go channel
func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	c, err := q.channel.Consume(q.Name,
		"", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	return c
}

// Close关闭消息队列
func (q *RabbitMQ) Close() {
	err := q.channel.Close()
	if err != nil {
		log.Println(err)
	}
	err = q.conn.Close()
	if err != nil {
		log.Println(err)
	}
}
