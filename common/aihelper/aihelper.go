package aihelper

import (
	"GoNexus/common/rabbitmq"
	"GoNexus/model"
	"sync"
)

// AIHelper AI助手结构体
type AIHelper struct {
	model     AIModel                                      // AI模型接口,支持不同模型实现
	messages  []*model.Message                             // 历史消息列表,存储用户与AI之间的对话记录
	mu        sync.Mutex                                   // 读写锁,保护历史消息并发访问
	SessionID string                                       // 会话唯一标识,用于绑定消息和上下文
	saveFunc  func(*model.Message) (*model.Message, error) // 消息存储回调函数,异步发布到RabbitMQ
}

// NewAIHelper 创建新的AIHelper实例方法
func NewAIHelper(aiModel AIModel, SessionID string) *AIHelper {
	return &AIHelper{
		model:     aiModel,
		messages:  make([]*model.Message, 0),
		mu:        sync.Mutex{},
		SessionID: SessionID,
		// 异步推送到消息队列中
		saveFunc: func(msg *model.Message) (*model.Message, error) {
			data := rabbitmq.GenerateMessageMQParam(msg.SessionID, msg.UserName, msg.Content, msg.IsUser)
			err := rabbitmq.RMQMessage.Publish(data)
			return msg, err
		},
	}
}
