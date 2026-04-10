package tts

import (
	"GoNexus/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// BaiduTTS 百度智能云TTS服务结构体
type BaiduTTS struct {
}

// NewBaiduTTS 创建百度云TTS服务结构体实例方法
func NewBaiduTTS() *BaiduTTS {
	return &BaiduTTS{}
}

// BaiduTTSCreateTaskRequest 百度云TTS服务创建任务接口请求结构体
type BaiduTTSCreateTaskRequest struct {
	Text           string `json:"text"`
	Format         string `json:"format"`
	Voice          int    `json:"voice"`
	Lang           string `json:"lang"`
	Speed          int    `json:"speed"`
	Pitch          int    `json:"pitch"`
	Volume         int    `json:"volume"`
	EnableSubtitle int    `json:"enable_subtitle"`
}

// BaiduTTSCreateTaskResponse 百度云TTS创建任务接口响应结构体
type BaiduTTSCreateTaskResponse struct {
	TaskID string `json:"task_id"`
}

// CreateTTS 百度云TTS长文本在线合成的创建任务方法
func (t *BaiduTTS) CreateTTS(ctx context.Context, text string) (string, error) {
	accessToken := t.GetAccessToken()
	if accessToken == "" {
		return "", fmt.Errorf("get Baidu tts access token failed")
	}
	url := "https://aip.baidubce.com/rpc/2.0/tts/v1/create?access_token=" + accessToken
	payload := BaiduTTSCreateTaskRequest{
		Text:           text,
		Format:         "mp3-16k",
		Voice:          4194,
		Lang:           "zh",
		Speed:          5,
		Pitch:          5,
		Volume:         5,
		EnableSubtitle: 0,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	log.Println("[Baidu TTS Create] raw:", string(body))
	result := BaiduTTSCreateTaskResponse{}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.TaskID == "" {
		return "", fmt.Errorf("create Baidu tts failed, empty task_id")
	}
	return result.TaskID, nil
}

// GetAccessToken 获取TTS鉴权签名信息
func (t *BaiduTTS) GetAccessToken() string {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf(
		"grant_type=client_credentials&client_id=%s&client_secret=%s",
		config.GetConfig().VoiceServiceApiKey,
		config.GetConfig().VoiceServiceSecretKey)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		log.Printf("get token failed. err: %v", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read token failed. err: %v", err)
		return ""
	}
	accessTokenObj := map[string]any{}
	_ = json.Unmarshal(body, &accessTokenObj)
	return accessTokenObj["access_token"].(string)
}

type TTSTaskResult struct {
	SpeechURL string `json:"speech_url,omitempty"`
}

type TTSTask struct {
	TaskID     string         `json:"task_id"`
	TaskStatus string         `json:"task_status"`
	TaskResult *TTSTaskResult `json:"task_result,omitempty"`
}

type TTSQuryResponse struct {
	LogID     string    `json:"log_id"`
	TasksInfo []TTSTask `json:"tasks_info"`
}

// QueryTTSFull 百度云TTS长文本在线合成的查询任务结果
func (t *BaiduTTS) QueryTTSFull(ctx context.Context, taskID string) (*TTSQuryResponse, error) {
	accessToken := t.GetAccessToken()
	if accessToken == "" {
		return nil, fmt.Errorf("get access token failed")
	}
	reqBody := map[string][]string{
		"task_ids": {taskID},
	}
	jsonBody, _ := json.Marshal(reqBody)
	url := "https://aip.baidubce.com/rpc/2.0/tts/v1/query?access_token=" + accessToken
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println("[TTS Query] raw:", string(body))
	var rawRes struct {
		LogID     json.Number `json:"log_id"`
		TasksInfo []struct {
			TaskID     string          `json:"task_id"`
			TaskStatus string          `json:"task_status"`
			TaskResult json.RawMessage `json:"task_result,omitempty"`
		} `json:"tasks_info"`
	}
	if err = json.Unmarshal(body, &rawRes); err != nil {
		return nil, err
	}
	result := &TTSQuryResponse{
		LogID:     rawRes.LogID.String(),
		TasksInfo: make([]TTSTask, 0, len(rawRes.TasksInfo)),
	}
	for _, taskInfo := range rawRes.TasksInfo {
		task := TTSTask{
			TaskID:     taskInfo.TaskID,
			TaskStatus: taskInfo.TaskStatus,
			TaskResult: nil,
		}
		if taskInfo.TaskStatus == "Success" && len(taskInfo.TaskResult) > 0 {
			var r TTSTaskResult
			if err = json.Unmarshal(taskInfo.TaskResult, &r); err != nil {
				log.Println("parse task_result error:", err)
				return nil, fmt.Errorf("failed to parse task result: %v", err)
			}
			task.TaskResult = &r
		}
		result.TasksInfo = append(result.TasksInfo, task)
	}
	return result, nil
}
