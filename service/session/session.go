package session

import (
	"GoNexus/common/aihelper"
	"GoNexus/common/code"
	"GoNexus/dao/session"
	"GoNexus/model"
	"context"
	"log"
	"net/http"

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
		return "", "", code.AIModelFail
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
		return "", code.AIModelFail
	}
	// 2. 生成AI回复,同步对话
	aiResponse, err := helper.GenerateResponse(context.Background(), username, userQuestion)
	if err != nil {
		log.Println("ChatSend GenerateResponse failed. err:", err)
		return "", code.AIModelFail
	}
	return aiResponse.Content, code.SuccessCode
}

// CreateStreamSessionOnly 只创建流式会话
func CreateStreamSessionOnly(username, userQuestion string) (string, code.Code) {
	newSession := &model.Session{
		ID:       uuid.NewString(),
		Username: username,
		Title:    userQuestion,
	}
	createdSession, err := session.CreateSession(newSession)
	if err != nil {
		log.Println("CreateStreamSessionOnly CreateSession failed. err:", err)
		return "", code.ServerBusyCode
	}
	return createdSession.ID, code.SuccessCode
}

// StreamMessageToCurrentSession 以流式方式在当前会话中进行传输信息
func StreamMessageToCurrentSession(username, sessionID, userQuestion, modelType string, writer http.ResponseWriter) code.Code {
	// 1. 确保writer支持Flush
	flusher, ok := writer.(http.Flusher)
	if !ok {
		log.Println("StreamMessageToCurrentSession: streaming unsupported")
		return code.ServerBusyCode
	}
	// 2. 获取AIHelper实例
	globalManager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"username": username,
	}
	helper, err := globalManager.GetOrCreateAIHelper(username, sessionID, modelType, config)
	if err != nil {
		log.Println("StreamMessageToCurrentSession GetOrCreateAIHelper failed. err:", err)
		return code.AIModelFail
	}
	// 3. 定义回调方法
	cb := func(msg string) {
		log.Printf("[SSE] Sending chunk: %s (len=%d)\n", msg, len(msg))
		_, err = writer.Write([]byte("data: " + msg + "\n\n"))
		if err != nil {
			log.Println("[SSE] Write error:", err)
			return
		}
		flusher.Flush()
		log.Println("[SSE] Flushed")
	}
	// 4. 调用AI对话
	_, err = helper.StreamResponse(context.Background(), username, userQuestion, cb)
	if err != nil {
		log.Println("StreamMessageToCurrentSession StreamResponse failed. err:", err)
		return code.AIModelFail
	}
	// 5. 传回前端完成标记
	_, err = writer.Write([]byte("data: [DONE]\n\n"))
	if err != nil {
		log.Println("[SSE] Write DONE error:", err)
		return code.AIModelFail
	}
	flusher.Flush()
	return code.SuccessCode
}

// ChatStreamSend 基于当前会话窗口与AI流式聊天
func ChatStreamSend(username, sessionID, userQuestion, modelType string, writer http.ResponseWriter) code.Code {
	return StreamMessageToCurrentSession(username, sessionID, userQuestion, modelType, writer)
}
