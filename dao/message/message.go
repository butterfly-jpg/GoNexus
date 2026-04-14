package message

import (
	"GoNexus/common/mysql"
	"GoNexus/model"
)

// CreateMessage 在数据库Message表中创建一条消息数据
func CreateMessage(message *model.Message) (*model.Message, error) {
	err := mysql.DB.Create(message).Error
	return message, err
}

// DeleteMessagesBySessionID 根据会话ID删除该会话下的所有消息
func DeleteMessagesBySessionID(sessionID string) error {
	return mysql.DB.Where("session_id = ?", sessionID).Delete(&model.Message{}).Error
}

// GetAllMessages 从数据库Message表中查询所有的数据
func GetAllMessages() ([]*model.Message, error) {
	var messages []*model.Message
	err := mysql.DB.Order("created_at asc").Find(&messages).Error
	return messages, err
}
