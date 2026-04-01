package aihelper

import (
	"GoNexus/common/rag"
	"GoNexus/config"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
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
	apiKey := config.GetConfig().DeepSeekApiKey
	if apiKey == "" {
		return nil, errors.New("deepseekApiKey config is not set")
	}
	llm, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  apiKey,
		BaseURL: config.GetConfig().DeepSeekBaseUrl,
		Model:   config.GetConfig().DeepSeekModelName,
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
	apiKey := config.GetConfig().QwenApiKey
	if apiKey == "" {
		return nil, errors.New("qwenApiKey config is not set")
	}
	llm, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:  apiKey,
		BaseURL: config.GetConfig().QwenBaseUrl,
		Model:   config.GetConfig().QwenModelName,
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

// =========================== RAG 实现(基于千问模型) =========================

// QwenRAGModel 接入Qwen大模型的Rag服务结构体
type QwenRAGModel struct {
	llm      model.ToolCallingChatModel
	username string
}

// NewQwenRAGModel 获取Qwen-Rag模型实例
func NewQwenRAGModel(ctx context.Context, username string) (*QwenRAGModel, error) {
	apiKey := config.GetConfig().QwenApiKey
	if apiKey == "" {
		return nil, errors.New("qwenApiKey config is not set")
	}
	llm, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:  apiKey,
		BaseURL: config.GetConfig().QwenBaseUrl,
		Model:   config.GetConfig().QwenModelName,
	})
	if err != nil {
		return nil, fmt.Errorf("create qwen rag model failed. err: %v", err)
	}
	return &QwenRAGModel{llm: llm, username: username}, nil
}

// GenerateResponse QwenRag同步生成回复方法实现
func (q *QwenRAGModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	// 1. 创建RAG查询器
	ragQuery, err := rag.NewRAGQuery(ctx, q.username)
	if err != nil {
		log.Printf("create rag query failed(user may not have uploaded file). err: %v", err)
		// 用户没有上传文件，那么正常调用Qwen对用户问题进行回答
		res, err := q.llm.Generate(ctx, messages)
		if err != nil {
			return nil, fmt.Errorf("qwen rag generate failed. err: %v", err)
		}
		return res, nil
	}
	// 2. 获取用户的最后一条消息提问
	if len(messages) == 0 {
		return nil, fmt.Errorf("user message is empty")
	}
	lastMessage := messages[len(messages)-1]
	// 3. 检索相关文档
	docs, err := ragQuery.RetrieveDocuments(ctx, lastMessage.Content)
	if err != nil {
		log.Printf("retrieve documents failed. err: %v", err)
		// 检索文档失败，还是正常调用Qwen对用户问题进行回答
		res, err := q.llm.Generate(ctx, messages)
		if err != nil {
			return nil, fmt.Errorf("qwen rag generate failed. err: %v", err)
		}
		return res, nil
	}
	// 4. 将用户的提问和检索后的结果组装在一起构建新的RAG提示词
	ragPrompt := rag.BuildRAGPrompt(lastMessage.Content, docs)
	// 5. 将最后一条消息替换为新的RAG提示词
	ragMessages := make([]*schema.Message, len(messages))
	copy(ragMessages, messages)
	ragMessages[len(ragMessages)-1] = &schema.Message{
		Role:    schema.User,
		Content: ragPrompt,
	}
	// 6. 调用LLM生成回答
	res, err := q.llm.Generate(ctx, ragMessages)
	if err != nil {
		return nil, fmt.Errorf("qwen rag generate failed. err: %v", err)
	}
	return res, nil
}

// StreamResponse QwenRag流式生成回复方法实现
func (q *QwenRAGModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	// 1. 创建RAG查询器
	ragQuery, err := rag.NewRAGQuery(ctx, q.username)
	if err != nil {
		log.Printf("create rag query failed(user may not have uploaded file). err: %v", err)
		// 用户没有上传文件，那么正常调用Qwen对用户问题进行流式回答
		return q.StreamWithoutRAG(ctx, messages, cb)
	}
	// 2. 获取用户的最后一条消息提问
	if len(messages) == 0 {
		return "", fmt.Errorf("user message is empty")
	}
	lastMessage := messages[len(messages)-1]
	// 3. 检索相关文档
	docs, err := ragQuery.RetrieveDocuments(ctx, lastMessage.Content)
	if err != nil {
		log.Printf("retrieve documents failed. err: %v", err)
		return q.StreamWithoutRAG(ctx, messages, cb)
	}
	// 4. 用户提问和检索得到的结果组装替换最后一条消息
	ragPrompt := rag.BuildRAGPrompt(lastMessage.Content, docs)
	ragMessages := make([]*schema.Message, len(messages))
	copy(ragMessages, messages)
	ragMessages[len(ragMessages)-1] = &schema.Message{
		Role:    schema.User,
		Content: ragPrompt,
	}
	// 5. 调用LLM流式生成回答
	stream, err := q.llm.Stream(ctx, ragMessages)
	if err != nil {
		return "", fmt.Errorf("qwen rag stream failed. err: %v", err)
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

// StreamWithoutRAG 没有RAG文档时的原始流式输出
func (q *QwenRAGModel) StreamWithoutRAG(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
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

// GetModelType 获取模型类型,QwenRag是3号模型
func (q *QwenRAGModel) GetModelType() string {
	return "3"
}
