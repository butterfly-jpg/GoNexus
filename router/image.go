package router

import (
	"GoNexus/controller/image"

	"github.com/gin-gonic/gin"
)

// RegisterImageRouter 图像识别接口路由注册
func RegisterImageRouter(r *gin.RouterGroup) {
	r.POST("/recognize", image.RecognizeImage)
}
