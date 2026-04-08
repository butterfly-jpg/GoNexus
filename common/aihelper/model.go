package aihelper

import (
	"GoNexus/common/rag"
	"GoNexus/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
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
		return q.StreamByQwenRag(ctx, messages, cb)
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
		return q.StreamByQwenRag(ctx, messages, cb)
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
	return q.StreamByQwenRag(ctx, ragMessages, cb)
}

// StreamByQwenRag 基于QwenRag的原始流式输出
func (q *QwenRAGModel) StreamByQwenRag(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
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

// =========================== MCP 实现(基于千问模型) =========================

// QwenMCPModel 基于Qwen的MCP模型实现,集成MCP服务
type QwenMCPModel struct {
	llm        model.ToolCallingChatModel
	mcpClient  *client.Client
	username   string
	mcpBaseURL string
}

// NewQwenMCPModel 创建QwenMCP模型实例
func NewQwenMCPModel(ctx context.Context, username string) (*QwenMCPModel, error) {
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
	mcpBaseURL := "http://localhost:8081/mcp"
	return &QwenMCPModel{
		llm:        llm,
		username:   username,
		mcpBaseURL: mcpBaseURL,
	}, nil
}

// GenerateResponse QwenMCP同步生成消息回复方法实现
func (m *QwenMCPModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("messages is empty")
	}
	// 1. 获取消息的最后一条
	lastMessage := messages[len(messages)-1]
	query := lastMessage.Content
	// 2. 第一次调用AI
	// 2.1 将用户的问题第一次转为AI擅长的Prompt格式
	firstPrompt := m.buildFirstPrompt(query)
	// 2.2 得到第一次合并整理的消息
	firstMessages := make([]*schema.Message, len(messages))
	copy(firstMessages, messages)
	firstMessages[len(firstMessages)-1] = &schema.Message{
		Role:    schema.User,
		Content: firstPrompt,
	}
	// 2.3 调用LLM生成第一次响应
	firstRes, err := m.llm.Generate(ctx, firstMessages)
	if err != nil {
		return nil, fmt.Errorf("qwen mcp first generate failed. err: %v", err)
	}
	// 3. 解析AI响应
	aiResult := firstRes.Content
	toolCall, err := m.parseAIResponse(aiResult)
	if err != nil {
		log.Printf("parseAIResponse failed. err: %v", err)
		return firstRes, nil
	}
	// 4. 判断是否需要调用MCP工具
	// 4.1 AI不调用工具,直接返回响应
	if !toolCall.IsToolCall {
		log.Println("toolCall IsToolCall is false")
		return firstRes, nil
	}
	// 4.2 AI需要调用工具
	log.Println("toolCall IsToolCall is true")
	// 4.2.1 获取MCP客户端
	mcpClient, err := m.getMCPClient(ctx)
	if err != nil {
		log.Printf("getMCPClient failed. err: %v", err)
		return firstRes, nil
	}
	// 4.2.2 调用MCP工具
	toolResult, err := m.callMCPTool(ctx, mcpClient, toolCall.ToolName, toolCall.Args)
	if err != nil {
		log.Printf("callMCPTool failed. err: %v", err)
		return firstRes, nil
	}
	// 5. 第二次调用AI
	// 5.1 将用户的问题第二次转为AI擅长的Prompt格式
	secondPrompt := m.buildSecondPrompt(query, toolCall.ToolName, toolResult, toolCall.Args)
	// 5.2 得到第二次合并整理的消息
	secondMessages := make([]*schema.Message, len(firstMessages))
	copy(secondMessages, firstMessages)
	secondMessages[len(secondMessages)-1] = &schema.Message{
		Role:    schema.User,
		Content: secondPrompt,
	}
	// 5.3 调用LLM生成第二次响应为最终响应
	finalRes, err := m.llm.Generate(ctx, secondMessages)
	if err != nil {
		return nil, fmt.Errorf("qwen mcp second generate failed. err: %v", err)
	}
	log.Println("最终响应为：", finalRes)
	return finalRes, nil
}

// StreamResponse QwenMCP流式生成消息回复方法实现
func (m *QwenMCPModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("messages is empty")
	}
	// 1. 获取消息的最后一条
	lastMessage := messages[len(messages)-1]
	query := lastMessage.Content
	// 2. 第一次调用AI
	// 2.1 将用户的问题第一次转为AI擅长的Prompt格式
	firstPrompt := m.buildFirstPrompt(query)
	// 2.2 得到第一次合并整理的消息
	firstMessages := make([]*schema.Message, len(messages))
	copy(firstMessages, messages)
	firstMessages[len(firstMessages)-1] = &schema.Message{Role: schema.User, Content: firstPrompt}
	// 2.3 调用LLM生成第一次响应
	firstRes, err := m.llm.Generate(ctx, firstMessages)
	if err != nil {
		return "", fmt.Errorf("qwen mcp first generate failed. err: %v", err)
	}
	// 3. 解析AI响应
	toolCall, err := m.parseAIResponse(firstRes.Content)
	// 4. 判断是否需要调用MCP工具
	// 4.1 AI不调用工具,直接返回响应
	if err != nil || !toolCall.IsToolCall {
		stream, err := m.llm.Stream(ctx, firstMessages)
		if err != nil {
			return "", fmt.Errorf("qwen mcp stream failed. err: %v", err)
		}
		defer stream.Close()
		var fullRes strings.Builder
		for {
			msg, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return "", fmt.Errorf("qwen mcp stream receive failed. err: %v", err)
			}
			if len(msg.Content) > 0 {
				fullRes.WriteString(msg.Content)
				cb(msg.Content)
			}
		}
		return fullRes.String(), nil
	}
	// 4.2 AI需要调用工具
	// 4.2.1 获取MCP客户端
	mcpClient, err := m.getMCPClient(ctx)
	if err != nil {
		return "", fmt.Errorf("getMCPClient failed. err: %v", err)
	}
	// 4.2.2 调用MCP工具
	toolResult, err := m.callMCPTool(ctx, mcpClient, toolCall.ToolName, toolCall.Args)
	if err != nil {
		return "", fmt.Errorf("callMCPTool failed. err: %v", err)
	}
	// 5. 第二次调用AI
	// 5.1 将用户的问题第二次转为AI擅长的Prompt格式
	secondPrompt := m.buildSecondPrompt(query, toolCall.ToolName, toolResult, toolCall.Args)
	// 5.2 得到第二次合并整理的消息
	secondMessages := make([]*schema.Message, len(firstMessages))
	copy(secondMessages, firstMessages)
	secondMessages[len(secondMessages)-1] = &schema.Message{Role: schema.User, Content: secondPrompt}
	// 5.3 调用LLM生成第二次响应为最终响应
	stream, err := m.llm.Stream(ctx, secondMessages)
	if err != nil {
		return "", fmt.Errorf("qwen mcp second stream failed. err: %v", err)
	}
	defer stream.Close()
	var fullRes strings.Builder
	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("qwen mcp second stream receive failed. err: %v", err)
		}
		if len(msg.Content) > 0 {
			fullRes.WriteString(msg.Content)
			cb(msg.Content)
		}
	}
	return fullRes.String(), nil
}

// GetModelType 获取模型类型,QwenMCP是4号模型
func (m *QwenMCPModel) GetModelType() string {
	return "4"
}

// buildFirstPrompt 构建第一次调用的提示词
func (m *QwenMCPModel) buildFirstPrompt(query string) string {
	return fmt.Sprintf(`你是一个智能助手，可以调用MCP工具来获取信息。

可用工具:
- get_weather: 获取指定城市的天气信息，参数: city（城市名称，支持中文和英文，如北京、Shanghai等）

重要规则:
1. 如果需要调用工具，必须严格返回以下JSON格式:
{
  "isToolCall": true,
  "toolName": "工具名称",
  "args": {"参数名": "参数值"}
}
2. 如果不需要调用工具，直接返回自然语言回答
3. 请根据用户问题决定是否需要调用工具

用户问题: %s

请根据需要调用适当的工具，然后给出综合的回答。`, query)
}

// buildSecondPrompt 构建第二次调用的提示词
func (m *QwenMCPModel) buildSecondPrompt(query, toolName, toolResult string, args map[string]interface{}) string {
	return fmt.Sprintf(`你是一个智能助手，可以调用MCP工具来获取信息。

工具执行结果:
工具名称: %s
工具参数: %v
工具结果: %s

用户问题: %s

请根据工具结果和用户问题，给出最终的综合回答。`, toolName, args, toolResult, query)
}

// AIToolCall 表示AI工具调用请求
type AIToolCall struct {
	IsToolCall bool                   `json:"isToolCall"`
	ToolName   string                 `json:"toolName"`
	Args       map[string]interface{} `json:"args"`
}

// parseAIResponse 解析AI响应,检查是否包含工具调用
func (m *QwenMCPModel) parseAIResponse(response string) (*AIToolCall, error) {
	// 1. 尝试解析为JSON
	// 1.1 解析JSON成功,直接返回
	var aiToolCall AIToolCall
	if err := json.Unmarshal([]byte(response), &aiToolCall); err == nil {
		return &aiToolCall, nil
	}
	// 1.2 解析JSON失败,检查是否包含工具调用关键词
	if strings.Contains(response, "get_weather") {
		// 包含但JSON无法解析,city使用空字符串
		return &AIToolCall{
			IsToolCall: true,
			ToolName:   "get_weather",
			Args:       map[string]interface{}{"city": ""},
		}, nil
	}
	// 1.3 不包含工具调用关键词,说明不需要调用MCP工具
	return &AIToolCall{
		IsToolCall: false,
	}, nil
}

// getMCPClient 获取或创建MCP客户端
func (m *QwenMCPModel) getMCPClient(ctx context.Context) (*client.Client, error) {
	if m.mcpClient == nil {
		// 1. 创建MCP客户端
		httpTransport, err := transport.NewStreamableHTTP(m.mcpBaseURL)
		if err != nil {
			return nil, fmt.Errorf("create mcp transport failed. err: %v", err)
		}
		m.mcpClient = client.NewClient(httpTransport)
		// 2. 初始化MCP客户端
		initializeRequest := mcp.InitializeRequest{}
		initializeRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initializeRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "MCP-Go Weather Client",
			Version: "1.0.0",
		}
		initializeRequest.Params.Capabilities = mcp.ClientCapabilities{}
		if _, err = m.mcpClient.Initialize(ctx, initializeRequest); err != nil {
			return nil, fmt.Errorf("mcp client init failed. err: %v", err)
		}
	}
	return m.mcpClient, nil
}

// callMCPTool 调用MCP工具
func (m *QwenMCPModel) callMCPTool(ctx context.Context, client *client.Client, toolName string, args map[string]interface{}) (string, error) {
	// 1. 构建MCP请求
	callToolRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: args,
		},
	}
	// 2. 发送请求调用MCP
	result, err := client.CallTool(ctx, callToolRequest)
	if err != nil {
		return "", fmt.Errorf("call mcp tool failed. err: %v", err)
	}
	// 3. 解析结果提取文本
	var text string
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			text += textContent.Text + "\n"
		}
	}
	return text, nil
}
