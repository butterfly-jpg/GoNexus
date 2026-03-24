package session

import (
	"GoNexus/common/code"
	"GoNexus/controller"
	"GoNexus/service/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	CreateSessionAndSendMessageRequest struct {
		UserQuestion string `json:"question" binding:"required"`
		ModelType    string `json:"modelType" binding:"required"`
	}
	CreateSessionAndSendMessageResponse struct {
		AiInformation string `json:"information,omitempty"`
		SessionID     string `json:"sessionID,omitempty"`
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
