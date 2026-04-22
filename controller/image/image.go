package image

import (
	"GoNexus/common/code"
	"GoNexus/controller"
	"GoNexus/service/image"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	// RecognizeImageResponse 识别图片接口方法响应结构
	RecognizeImageResponse struct {
		ClassName string `json:"class_name,omitempty"` // AI回答
		controller.Response
	}
)

// RecognizeImage 控制层识别图片接口
func RecognizeImage(c *gin.Context) {
	res := new(RecognizeImageResponse)
	file, err := c.FormFile("image")
	if err != nil {
		log.Println("FormFile err:", err)
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	className, err := image.RecognizeImage(file)
	if err != nil {
		log.Println("RecognizeImage err:", err)
		c.JSON(http.StatusOK, res.CodeOf(code.ServerBusyCode))
		return
	}
	res.Success()
	res.ClassName = className
	c.JSON(http.StatusOK, res)
}
