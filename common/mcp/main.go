package main

import (
	mcpclient "GoNexus/common/mcp/client"
	mcpserver "GoNexus/common/mcp/server"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// 1. 定义命令行参数标志
	mode := flag.String("mode", "", "运行模式: server 或 client")
	httpAddr := flag.String("http-addr", ":8082", "HTTP服务器地址")
	city := flag.String("city", "", "要查询天气的城市名称")
	flag.Parse()
	// 2. 校验参数
	if *mode == "" {
		log.Println("Error: 您必须指定模式使用--mode (server 或 client)")
		flag.Usage()
		os.Exit(1)
	}
	// 3. 启动MCP服务器或MCP客户端
	if *mode == "server" {
		// 3.1 启动服务器
		log.Println("MCP server starting...")
		if err := mcpserver.StartServer(*httpAddr); err != nil {
			log.Fatalf("MCP server start failed: %v", err)
		}
	} else if *mode == "client" {
		// 3.2 启动客户端
		// a.校验参数
		if *city == "" {
			fmt.Println("Error: 您必须指定城市名称使用--city")
			flag.Usage()
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// b.创建客户端
		httpURL := "http://localhost:8082/mcp"
		mcpClient, err := mcpclient.NewMCPClient(httpURL)
		if err != nil {
			log.Fatalf("create MCP client failed: %v", err)
		}
		defer mcpClient.Close()
		// c.初始化客户端
		if _, err = mcpClient.Initialize(ctx); err != nil {
			log.Fatalf("initialize MCP client failed: %v", err)
		}
		// d.执行健康检查
		if err = mcpClient.Ping(ctx); err != nil {
			log.Fatalf("health check failed: %v", err)
		}
		// e.调用天气工具
		result, err := mcpClient.CallWeatherTool(ctx, *city)
		if err != nil {
			log.Fatalf("call weather tool failed: %v", err)
		}
		// f.显示天气结果
		fmt.Println("\n天气查询结果:")
		fmt.Println(mcpClient.GetToolResultText(result))
		fmt.Println("\n客户端初始化成功，正在关闭...")
	}
}
