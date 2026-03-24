package session

import (
	"GoNexus/common/aihelper"
	"GoNexus/common/code"
	"GoNexus/dao/session"
	"GoNexus/model"
	"log"

	"github.com/google/uuid"
)

// CreateSessionAndSendMessage 创建会话和发送消息方法
func CreateSessionAndSendMessage(userName, userQuestion, modelType string) (string, string, code.Code) {
	// 1. 创建新会话
	newSession := &model.Session{
		ID:       uuid.NewString(),
		UserName: userName,
		Title:    userQuestion, // 使用用户第一次的问题作为会话标题
	}
	createdSession, err := session.CreateSession(newSession)
	if err != nil {
		log.Println("CreateSessionAndSendMessage CreateSession failed. err:", err)
		return "", "", code.ServerBusyCode
	}
	// 2. 创建 AIHelper 实例
	globalManager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"apiKey":   "api-key", // todo
		"username": userName,
	}
	helper, err := globalManager.GetOrCreateAIHelper(userName, createdSession.ID, modelType, config)
	if err != nil {
		log.Println("CreateSessionAndSendMessage GetOrCreateAIHelper failed. err:", err)
		return "", "", code.ServerBusyCode
	}
	// 3. 生成AI回复

	return "", "", code.SuccessCode
}
