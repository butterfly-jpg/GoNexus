package session

import (
	"GoNexus/common/mysql"
	"GoNexus/model"
)

// CreateSession 在数据库Session表中创建一条会话数据
func CreateSession(session *model.Session) (*model.Session, error) {
	err := mysql.DB.Create(session).Error
	return session, err
}
