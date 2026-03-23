package rabbitmq

import (
	"GoNexus/config"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// ExchangeName 自定义exchange名字
	ExchangeName = "GoNexus"
	// ExchangeType 自定义exchange类型
	ExchangeType = "direct"
	// RoutingKey 自定义路由键
	RoutingKey = "Message"
)

// conn 全局连接对象
var conn *amqp.Connection

// initConn 初始化connection对象
func initConn() {
	c := config.GetConfig()
	mqUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.RabbitmqUsername, c.RabbitmqPassword, c.RabbitmqHost, c.RabbitmqPort, c.RabbitmqVhost)
	var err error
	conn, err = amqp.Dial(mqUrl)
	if err != nil {
		log.Fatalf("connect to RabbitMQ failed. err: %v", err)
	}
}

// RabbitMQ RabbitMQ结构体
type RabbitMQ struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	Exchange   string
	routingKey string
}

// NewRabbitMQ 创建自定义Exchange的RabbitMQ对象
func NewRabbitMQ(exchangeName, exchangeType, routingKey string) *RabbitMQ {
	// 1. 获取connection
	if conn == nil {
		initConn()
	}
	// 2. 创建channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("create channel failed. err: %v", err)
	}
	// 3. 声明Exchange
	if err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatalf("declare exchange failed. err: %v", err)
	}
	// 4. 返回实例
	return &RabbitMQ{
		conn:       conn,
		channel:    ch,
		Exchange:   exchangeName,
		routingKey: routingKey,
	}
}

// Publish 生产者发送消息
func (r *RabbitMQ) Publish(message []byte) error {
	if err := r.channel.Publish(
		r.Exchange,
		r.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         message,
		},
	); err != nil {
		return fmt.Errorf("publish message failed. err: %v", err)
	}
	return nil
}

// Consume 消费者接收消息,handle：消息的消费业务函数,用于消费消息
func (r *RabbitMQ) Consume(handle func(msg *amqp.Delivery) error) {
	// 1. 创建/声明队列
	q, err := r.channel.QueueDeclare(
		r.routingKey, // 队列名==路由键
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("declare queue failed. err: %v", err)
	}
	// 2. 自定义交换机执行绑定
	if r.Exchange != "" {
		err = r.channel.QueueBind(
			q.Name,
			r.routingKey,
			r.Exchange,
			false,
			nil,
		)
		if err != nil {
			log.Printf("bind queue failed. err: %v", err)
		}
	}
	log.Printf("Queue %s bound to exchange %s with key %s", q.Name, r.Exchange, r.routingKey)
	// 3. 接收消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("consume message failed. err: %v", err)
	}
	// 4. 处理消息
	for msg := range msgs {
		if err = handle(&msg); err != nil {
			fmt.Printf("handle message failed. err: %v", err)
		}
	}
}
