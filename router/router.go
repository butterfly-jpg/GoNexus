package router

import (
	"GoNexus/middleware/jwt"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	enterRouter := r.Group("api/v1")
	// 1. 注册用户模块
	RegisterUserRouter(enterRouter.Group("/user"))
	// 2. 注册AI模块
	AIGroup := enterRouter.Group("/ai")
	AIGroup.Use(jwt.Auth())
	RegisterAIRouter(AIGroup)
	// 3. 注册文件上传模块(rag)
	FileGroup := enterRouter.Group("/file")
	FileGroup.Use(jwt.Auth())
	RegisterFileRouter(FileGroup)
	// 4. 注册图像识别模块
	ImageGroup := enterRouter.Group("/image")
	ImageGroup.Use(jwt.Auth())
	RegisterImageRouter(ImageGroup)
	return r
}
