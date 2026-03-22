package jwt

import (
	"GoNexus/common/code"
	"GoNexus/controller"
	"GoNexus/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth 拦截请求解析jwt token
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取请求头中 Authorization 字段值或 URL 中的token参数
		// Bearer格式 Authorization: Bearer <token>
		// URL格式 ?token=<token>
		var token string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = c.Query("token")
		}
		// 2. 验证token是否有效
		res := &controller.Response{}
		if token == "" {
			c.JSON(http.StatusOK, res.CodeOf(code.InvalidTokenCode))
			c.Abort()
			return
		}
		log.Println("token is:", token)
		userName, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, res.CodeOf(code.InvalidTokenCode))
			c.Abort()
			return
		}
		// 3. 将用户名存储到gin上下文中传递给后续业务逻辑
		c.Set("username", userName)
		c.Next()
	}
}
