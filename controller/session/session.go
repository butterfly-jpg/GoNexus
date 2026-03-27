package session

import (
	"GoNexus/common/code"
	"GoNexus/controller"
	"GoNexus/service/session"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	// CreateSessionAndSendMessageRequest 同步回复请求结构
	CreateSessionAndSendMessageRequest struct {
		UserQuestion string `json:"question" binding:"required"`
		ModelType    string `json:"modelType" binding:"required"`
	}
	// CreateSessionAndSendMessageResponse 同步回复响应结构
	CreateSessionAndSendMessageResponse struct {
		AiInformation string `json:"information,omitempty"`
		SessionID     string `json:"sessionID,omitempty"`
		controller.Response
	}
	// ChatSendRequest 同步聊天请求结构
	ChatSendRequest struct {
		UserQuestion string `json:"question" binding:"required"`
		ModelType    string `json:"modelType" binding:"required"`
		SessionID    string `json:"sessionID" binding:"required"`
	}
	// ChatSendResponse 同步聊天响应结构
	ChatSendResponse struct {
		AiInformation string `json:"information,omitempty"`
		controller.Response
	}
)

// CreateSessionAndSendMessage 控制层创建会话和发送消息方法
func CreateSessionAndSendMessage(c *gin.Context) {
	req := &CreateSessionAndSendMessageRequest{}
	res := &CreateSessionAndSendMessageResponse{}
	username := c.GetString("username")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	// 调用service层创建会话并发送消息的方法
	sessionId, aiInformation, resCode := session.CreateSessionAndSendMessage(username, req.UserQuestion, req.ModelType)
	if resCode != code.SuccessCode {
		c.JSON(http.StatusOK, res.CodeOf(resCode))
		return
	}
	res.Success()
	res.SessionID = sessionId
	res.AiInformation = aiInformation
	c.JSON(http.StatusOK, res)
}

// ChatSend 基于当前会话窗口与AI同步聊天
func ChatSend(c *gin.Context) {
	req := &ChatSendRequest{}
	res := &ChatSendResponse{}
	// 1. 获取JWT拦截器解析token获取的username
	username := c.GetString("username")
	// 2. 解析参数
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	// 3. 发送消息,将AI回复同步返回
	aiInformation, resCode := session.ChatSend(username, req.SessionID, req.UserQuestion, req.ModelType)
	if resCode != code.SuccessCode {
		c.JSON(http.StatusOK, res.CodeOf(resCode))
		return
	}
	res.Success()
	res.AiInformation = aiInformation
	c.JSON(http.StatusOK, res)
}

// CrateStreamSessionAndSendMessage 控制层创建会话和发送SSE消息回复方法
func CrateStreamSessionAndSendMessage(c *gin.Context) {
	// 1. 参数处理
	req := &CreateSessionAndSendMessageRequest{}
	username := c.GetString("username")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid parameters"})
		return
	}
	// 2. 设置SSE头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Accel-Buffering", "no")
	// 3. 先创建会话并获取sessionID
	sessionID, resCode := session.CreateStreamSessionOnly(username, req.UserQuestion)
	if resCode != code.SuccessCode {
		c.JSON(http.StatusOK, gin.H{"error": resCode})
		return
	}
	// 4. 将sessionID通过data事件发送给前端，使前端界面可以绑定当前会话，在侧边栏显示会话新标签
	c.Writer.WriteString(fmt.Sprintf("data: {\"sessionID\": \"%s\"}\n\n", sessionID))
	c.Writer.Flush()
	// 5. 将AI回复以流式模式发送给前端
	resCode = session.StreamMessageToCurrentSession(username, sessionID, req.UserQuestion, req.ModelType, http.ResponseWriter(c.Writer))
	if resCode != code.SuccessCode {
		c.SSEvent("error", gin.H{"message": "send message failed"})
		return
	}
}
