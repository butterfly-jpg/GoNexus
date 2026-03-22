package router

import "github.com/gin-gonic/gin"

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	// 1. 注册用户模块
	enterRouter := r.Group("api/v1")
	RegisterUserRouter(enterRouter.Group("/user"))
	// 2. 注册AI模块 todo

	// 3. 注册图像识别模块 todo

	// 4. 注册文件上传模块 todo
	return r
}
