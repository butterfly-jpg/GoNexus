package aihelper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type StreamCallback func(msg string)

// AIModel 定义AI模型接口
type AIModel interface {
	// GenerateResponse 同步生成回复
	GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error)
	// StreamResponse 流式生成回复，通过回调函数实时输出
	StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error)
	// GetModelType 返回模型类型
	GetModelType() string
}

// =========================== DeepSeek 实现 =========================

// DeepSeekModel DeepSeek模型实现
type DeepSeekModel struct {
	llm model.ToolCallingChatModel
}

// NewDeepSeekModel 获取DeepSeek模型实例
func NewDeepSeekModel(ctx context.Context) (*DeepSeekModel, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	modelName := os.Getenv("DEEPSEEK_MODEL_NAME")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	if apiKey == "" {
		return nil, errors.New("DEEPSEEK_API_KEY environment variable is not set")
	}
	llm, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   modelName,
	})
	if err != nil {
		return nil, fmt.Errorf("create deepseek model failed. err: %v", err)
	}
	return &DeepSeekModel{llm: llm}, nil
}

// GenerateResponse DeepSeek同步生成回复方法实现
func (ds *DeepSeekModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	res, err := ds.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("deepseek generate failed. err: %v", err)
	}
	return res, nil
}

// StreamResponse DeepSeek流式生成回复方法实现
func (ds *DeepSeekModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	stream, err := ds.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("deepseek stream failed. err: %v", err)
	}
	defer stream.Close()
	var fullRes strings.Builder
	for {
		// 从消息流中一帧一帧读取数据
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("deepseek stream receive failed. err: %v", err)
		}
		if len(msg.Content) > 0 {
			// 聚合AI回复的内容
			fullRes.WriteString(msg.Content)
			// 实时调用回调函数,将AI内容一帧一帧主动发送给前端
			cb(msg.Content)
		}
	}
	return fullRes.String(), nil
}

// GetModelType 获取模型类型,DeepSeek是1号模型
func (ds *DeepSeekModel) GetModelType() string {
	return "1"
}

// =========================== Qwen 实现 =========================

// QwenModel  Qwen模型实现
type QwenModel struct {
	llm model.ToolCallingChatModel
}

// NewQwenModel 获取Qwen模型实例
func NewQwenModel(ctx context.Context) (*QwenModel, error) {
	apiKey := os.Getenv("QWEN_API_KEY")
	modelName := os.Getenv("QWEN_MODEL_NAME")
	baseURL := os.Getenv("QWEN_BASE_URL")
	if apiKey == "" {
		return nil, errors.New("QWEN_API_KEY environment variable is not set")
	}
	llm, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   modelName,
	})
	if err != nil {
		return nil, fmt.Errorf("create qwen model failed. err: %v", err)
	}
	return &QwenModel{llm: llm}, nil
}

// GenerateResponse Qwen同步生成回复方法实现
func (q *QwenModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	res, err := q.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("qwen generate failed. err: %v", err)
	}
	return res, nil
}

// StreamResponse Qwen流式生成回复方法实现
func (q *QwenModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	stream, err := q.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("qwen stream failed. err: %v", err)
	}
	defer stream.Close()
	var fullRes strings.Builder
	for {
		// 从消息流中一帧一帧读取数据
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("qwen stream receive failed. err: %v", err)
		}
		if len(msg.Content) > 0 {
			// 聚合AI回复的内容
			fullRes.WriteString(msg.Content)
			// 实时调用回调函数,将AI内容一帧一帧主动发送给前端
			cb(msg.Content)
		}
	}
	return fullRes.String(), nil
}

// GetModelType 获取模型类型,Qwen是2号模型
func (q *QwenModel) GetModelType() string {
	return "2"
}
