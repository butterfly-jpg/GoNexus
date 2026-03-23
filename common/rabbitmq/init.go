package rabbitmq

var RMQMessage *RabbitMQ

// InitRabbitMQ 初始化RabbitMQ
func InitRabbitMQ() {
	// 创建MQ实例,不同队列共用一个connection,保证不同队列消费消息的顺序
	RMQMessage = NewRabbitMQ(ExchangeName, ExchangeType, RoutingKey)
	// 启动消费者
	go RMQMessage.Consume(MQMessage)
}
