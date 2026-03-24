package aihelper

import (
	"GoNexus/common/rabbitmq"
	"GoNexus/model"
	"GoNexus/utils"
	"context"
	"log"
	"sync"
)

// AIHelper AI助手结构体
type AIHelper struct {
	model     AIModel                                      // AI模型接口,支持不同模型实现
	messages  []*model.Message                             // 历史消息列表,存储用户与AI之间的对话记录
	mu        sync.RWMutex                                 // 读写锁,保护历史消息并发访问
	SessionID string                                       // 会话唯一标识,用于绑定消息和上下文
	saveFunc  func(*model.Message) (*model.Message, error) // 消息存储回调函数,异步发布到RabbitMQ
}

// NewAIHelper 创建新的AIHelper实例方法
func NewAIHelper(aiModel AIModel, SessionID string) *AIHelper {
	return &AIHelper{
		model:     aiModel,
		messages:  make([]*model.Message, 0),
		mu:        sync.RWMutex{},
		SessionID: SessionID,
		// 异步推送到消息队列中
		saveFunc: func(msg *model.Message) (*model.Message, error) {
			data := rabbitmq.GenerateMessageMQParam(msg.SessionID, msg.Username, msg.Content, msg.IsUser)
			err := rabbitmq.RMQMessage.Publish(data)
			return msg, err
		},
	}
}

// AddMessage 添加消息到内存中并异步存入RabbitMQ
func (h *AIHelper) AddMessage(content, username string, isUser, save bool) {
	userMsg := model.Message{
		SessionID: h.SessionID,
		Username:  username,
		Content:   content,
		IsUser:    isUser,
	}
	h.messages = append(h.messages, &userMsg)
	if save {
		_, err := h.saveFunc(&userMsg)
		if err != nil {
			log.Println("AddMessage save failed. err:", err)
		}
	}
}

// GenerateResponse 同步生成消息
func (h *AIHelper) GenerateResponse(ctx context.Context, username, userQuestion string) (*model.Message, error) {
	// 1. 存储用户消息model.Message
	h.AddMessage(userQuestion, username, true, true)
	// 2. 将model.Message转为schema.Message
	h.mu.RLock()
	messages := utils.ConvertToSchemaMessages(h.messages)
	h.mu.RUnlock()
	// 3. 调用大模型获得AI回复的消息schema.Message
	schemaMsg, err := h.model.GenerateResponse(ctx, messages)
	if err != nil {
		return nil, err
	}
	// 4. 将schema.Message转为model.Message
	modelMsg := utils.ConvertToModelMessages(h.SessionID, username, schemaMsg.Content)
	// 5. 存储AI消息schema.Message
	h.AddMessage(modelMsg.Content, username, false, true)
	return modelMsg, nil
}
