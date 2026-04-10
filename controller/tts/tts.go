package tts

import (
	"GoNexus/common/code"
	"GoNexus/controller"
	"GoNexus/service/tts"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	// TaskTTSRequest 创建TTS任务接口请求体
	TaskTTSRequest struct {
		Text string `json:"text,omitempty"`
	}
	// TaskTTSResponse 创建TTS任务接口响应体
	TaskTTSResponse struct {
		TaskID string `json:"task_id,omitempty"`
		controller.Response
	}
	// QueryTTSTaskResponse 查询TTS任务接口响应体
	QueryTTSTaskResponse struct {
		TaskID     string `json:"task_id,omitempty"`
		TaskStatus string `json:"task_status,omitempty"`
		TaskResult string `json:"task_result,omitempty"`
		controller.Response
	}
)

// CreateTTSTask 创建TTS任务
func CreateTTSTask(c *gin.Context) {
	req := new(TaskTTSRequest)
	res := new(TaskTTSResponse)
	baiduTTS := tts.NewBaiduTTSService()
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	if req.Text == "" {
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	taskID, err := baiduTTS.TTSService.CreateTTS(c, req.Text)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.TTSFail))
		return
	}
	res.Success()
	res.TaskID = taskID
	c.JSON(http.StatusOK, res)
}

// QueryTTSTask 查询TTS任务
func QueryTTSTask(c *gin.Context) {
	baiduTTS := tts.NewBaiduTTSService()
	taskID := c.Query("task_id")
	res := new(QueryTTSTaskResponse)
	if taskID == "" {
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	ttsQuryResponse, err := baiduTTS.TTSService.QueryTTSFull(c, taskID)
	if err != nil {
		log.Println("语音合成失败", err.Error())
		c.JSON(http.StatusOK, res.CodeOf(code.TTSFail))
		return
	}
	if len(ttsQuryResponse.TasksInfo) == 0 {
		c.JSON(http.StatusOK, res.CodeOf(code.TTSFail))
		return
	}
	res.Success()
	res.TaskID = ttsQuryResponse.TasksInfo[0].TaskID
	res.TaskStatus = ttsQuryResponse.TasksInfo[0].TaskStatus
	if ttsQuryResponse.TasksInfo[0].TaskResult != nil {
		res.TaskResult = ttsQuryResponse.TasksInfo[0].TaskResult.SpeechURL
	}
	c.JSON(http.StatusOK, res)
}
