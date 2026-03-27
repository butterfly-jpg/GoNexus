package session

import (
	"GoNexus/common/aihelper"
	"GoNexus/common/code"
	"GoNexus/dao/session"
	"GoNexus/model"
	"context"
	"log"

	"github.com/google/uuid"
)

// CreateSessionAndSendMessage 创建会话和发送消息方法,AI同步返回消息
func CreateSessionAndSendMessage(username, userQuestion, modelType string) (string, string, code.Code) {
	// 1. 创建新会话
	newSession := &model.Session{
		ID:       uuid.NewString(),
		Username: username,
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
		"username": username,
	}
	helper, err := globalManager.GetOrCreateAIHelper(username, createdSession.ID, modelType, config)
	if err != nil {
		log.Println("CreateSessionAndSendMessage GetOrCreateAIHelper failed. err:", err)
		return "", "", code.ServerBusyCode
	}
	// 3. 生成AI回复,同步对话
	aiResponse, err := helper.GenerateResponse(context.Background(), username, userQuestion)
	if err != nil {
		log.Println("CreateSessionAndSendMessage GenerateResponse failed. err:", err)
		return "", "", code.AIModelFail
	}
	return aiResponse.SessionID, aiResponse.Content, code.SuccessCode
}

// ChatSend 基于当前会话窗口与AI同步聊天
func ChatSend(username, sessionID, userQuestion, modelType string) (string, code.Code) {
	// 1. 获取AIHelper实例
	globalManager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"username": username,
	}
	helper, err := globalManager.GetOrCreateAIHelper(username, sessionID, modelType, config)
	if err != nil {
		log.Println("ChatSend GetOrCreateAIHelper failed. err:", err)
		return "", code.ServerBusyCode
	}
	// 2. 生成AI回复,同步对话
	aiResponse, err := helper.GenerateResponse(context.Background(), username, userQuestion)
	if err != nil {
		log.Println("ChatSend GenerateResponse failed. err:", err)
		return "", code.AIModelFail
	}
	return aiResponse.Content, code.SuccessCode
}
