package file

import (
	"GoNexus/common/code"
	"GoNexus/controller"
	"GoNexus/service/file"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	// UploadFileResponse 上传文件响应结构
	UploadFileResponse struct {
		FilePath string `json:"file_path,omitempty"`
		controller.Response
	}
)

// UploadingRagFile 上传文件构建RAG接口
func UploadingRagFile(c *gin.Context) {
	// 1. 参数处理
	res := &UploadFileResponse{}
	uploadedFile, err := c.FormFile("file")
	if err != nil {
		log.Println("FormFile err:", err)
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	username := c.GetString("username")
	// 2. 调用service层索引器
	filePath, err := file.UploadRagFile(username, uploadedFile)
	if err != nil {
		log.Println("UploadRagFile err:", err)
		c.JSON(http.StatusOK, res.CodeOf(code.ServerBusyCode))
		return
	}
	// 3. 返回
	res.Success()
	res.FilePath = filePath
	c.JSON(http.StatusOK, res)
}
