package router

import (
	"GoNexus/controller/session"
	"GoNexus/controller/tts"

	"github.com/gin-gonic/gin"
)

// RegisterAIRouter ai相关接口路由
func RegisterAIRouter(r *gin.RouterGroup) {
	// 展示用户会话列表
	r.GET("/chat/sessions", session.GetUserSessionByUsername)
	// 展示会话历史消息内容
	r.POST("/chat/history", session.ChatHistory)

	// 同步聊天接口
	r.POST("/chat/send-new-session", session.CreateSessionAndSendMessage)
	r.POST("/chat/send", session.ChatSend)

	// 流式聊天接口
	r.POST("/chat/send-stream-new-session", session.CreateStreamSessionAndSendMessage)
	r.POST("/chat/send-stream", session.ChatStreamSend)

	// TTS服务接口
	r.POST("/chat/tts", tts.CreateTTSTask)
	r.GET("/chat/tts/query", tts.QueryTTSTask)

	// 删除会话接口
	r.DELETE("/chat/delete-session", session.DeleteSession)
}
