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

// GetSessionBySessionId 根据SessionId查询对应的title数据
func GetSessionBySessionId(sessionId string) (*model.Session, error) {
	session := &model.Session{}
	err := mysql.DB.Where("id = ?", sessionId).First(session).Error
	return session, err
}

// GetSessionsByUsername 根据用户名查询该用户的所有会话
func GetSessionsByUsername(username string) ([]model.Session, error) {
	var sessions []model.Session
	err := mysql.DB.Where("username = ?", username).Order("created_at desc").Find(&sessions).Error
	return sessions, err
}

