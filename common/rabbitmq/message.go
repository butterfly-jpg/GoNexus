package rabbitmq

import (
	"GoNexus/dao/message"
	"GoNexus/model"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageMQParam MQ消息数据模型
type MessageMQParam struct {
	SessionID string `json:"session_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	IsUser    bool   `son:"is_user"`
}

// MQMessage 消费者将信息数据异步写入数据库
func MQMessage(msg *amqp.Delivery) error {
	var param MessageMQParam
	if err := json.Unmarshal(msg.Body, &param); err != nil {
		return err
	}
	newMsg := &model.Message{
		SessionID: param.SessionID,
		Username:  param.Username,
		Content:   param.Content,
		IsUser:    param.IsUser,
	}
	_, err := message.CreateMessage(newMsg)
	return err
}

// GenerateMessageMQParam 生成MQ消息数据方法
func GenerateMessageMQParam(sessionID, username, content string, isUer bool) []byte {
	param := MessageMQParam{
		SessionID: sessionID,
		Username:  username,
		Content:   content,
		IsUser:    isUer,
	}
	data, _ := json.Marshal(param)
	return data
}
