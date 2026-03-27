package router

import (
	"GoNexus/controller/session"

	"github.com/gin-gonic/gin"
)

// RegisterAIRouter ai相关接口路由
func RegisterAIRouter(r *gin.RouterGroup) {
	r.POST("/chat/send-new-session", session.CreateSessionAndSendMessage)
	r.POST("/chat/send", session.ChatSend)
}
