package router

import (
	"GoNexus/controller/file"

	"github.com/gin-gonic/gin"
)

// RegisterFileRouter 文件上传接口路由
func RegisterFileRouter(r *gin.RouterGroup) {
	r.POST("/upload", file.UploadingRagFile)
}
