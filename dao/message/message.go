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
