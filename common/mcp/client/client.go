package client

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCPClient MCP客户端的封装结构体
type MCPClient struct {
	c *client.Client
}

// NewMCPClient 创建新的MCP客户端实例
func NewMCPClient(httpURL string) (*MCPClient, error) {
	fmt.Println("正在初始化HTTP客户端...")
	httpTransport, err := transport.NewStreamableHTTP(httpURL)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP传输失败: %v", err)
	}
	c := client.NewClient(httpTransport)
	return &MCPClient{c}, nil
}

// Initialize 初始化MCP客户端
func (m *MCPClient) Initialize(ctx context.Context) (*mcp.InitializeResult, error) {
	// 1. 设置通知处理程序
	m.c.OnNotification(func(notification mcp.JSONRPCNotification) {
		fmt.Printf("收到通知: %s\n", notification.Method)
	})
	// 2. 初始化客户端
	fmt.Println("正在初始化客户端...")
	initializeRequest := mcp.InitializeRequest{}
	initializeRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initializeRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCP-Go Weather Client",
		Version: "1.0.0",
	}
	initializeRequest.Params.Capabilities = mcp.ClientCapabilities{}
	serverInfo, err := m.c.Initialize(ctx, initializeRequest)
	if err != nil {
		return nil, fmt.Errorf("初始化失败: %v", err)
	}
	// 3. 显示服务器信息
	fmt.Printf("连接到服务器: %s (版本 %s)\n", serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)
	return serverInfo, nil
}

// Close 关闭MCP客户端
func (m *MCPClient) Close() {
	if m.c != nil {
		m.c.Close()
	}
}

// Ping 执行健康检查
func (m *MCPClient) Ping(ctx context.Context) error {
	fmt.Println("正在执行健康检查...")
	if err := m.c.Ping(ctx); err != nil {
		return fmt.Errorf("健康检查失败: %v", err)
	}
	fmt.Println("服务器正常运行并响应")
	return nil
}

// CallWeatherTool 调用get_weather工具
func (m *MCPClient) CallWeatherTool(ctx context.Context, city string) (*mcp.CallToolResult, error) {
	fmt.Printf("正在查询城市 %s 的天气...\n", city)
	callToolRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "get_weather",
			Arguments: map[string]any{
				"city": city,
			},
		},
	}
	result, err := m.c.CallTool(ctx, callToolRequest)
	if err != nil {
		return nil, fmt.Errorf("调用工具失败: %v", err)
	}
	return result, nil
}

// GetToolResultText 获取工具结果中的文本内容
func (m *MCPClient) GetToolResultText(result *mcp.CallToolResult) string {
	var text string
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			text += textContent.Text + "\n"
		}
	}
	return text
}
