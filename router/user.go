package router

import (
	"GoNexus/controller/user"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.RouterGroup) {
	r.POST("/register", user.Register)
	r.POST("/captcha", user.HandleCaptcha)
	r.POST("/login", user.Login)
}
