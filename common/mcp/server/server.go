package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// WeatherAPIClient 天气API客户端结构体
type WeatherAPIClient struct {
}

// NewWeatherAPIClient 获取天气API客户端实例
func NewWeatherAPIClient() *WeatherAPIClient {
	return &WeatherAPIClient{}
}

// WeatherResponse 天气信息响应结构体
type WeatherResponse struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"windSpeed"`
}

// WttrResponse wtt.in的JSON响应
type WttrResponse struct {
	CurrentCondition []struct {
		TempC         string `json:"temp_C"`
		Humidity      string `json:"humidity"`
		WindspeedKmph string `json:"windspeedKmph"`
		WeatherDesc   []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
	} `json:"current_condition"`
	NearestArea []struct {
		AreaName []struct {
			Value string `json:"value"`
		} `json:"areaName"`
	} `json:"nearest_area"`
}

// GetWeather 获取天气信息方法
func (w *WeatherAPIClient) GetWeather(ctx context.Context, city string) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://wttr.in/%s?format=j1&lang=zh", city)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed. err: %v", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed. err: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed. err: %v", err)
	}
	wttrRes := &WttrResponse{}
	if err = json.Unmarshal(body, wttrRes); err != nil {
		return nil, fmt.Errorf("unmarshal response body failed. err: %v", err)
	}
	if len(wttrRes.CurrentCondition) == 0 {
		return nil, fmt.Errorf("no weather found for city %s", city)
	}
	cc := wttrRes.CurrentCondition[0]
	temp, _ := strconv.ParseFloat(cc.TempC, 64)
	humidity, _ := strconv.Atoi(cc.Humidity)
	wind, _ := strconv.ParseFloat(cc.WindspeedKmph, 64)
	location := city
	if len(wttrRes.NearestArea) > 0 && len(wttrRes.NearestArea[0].AreaName) > 0 {
		location = wttrRes.NearestArea[0].AreaName[0].Value
	}
	condition := ""
	if len(cc.WeatherDesc) > 0 {
		condition = cc.WeatherDesc[0].Value
	}
	return &WeatherResponse{
		Location:    location,
		Temperature: temp,
		Condition:   condition,
		Humidity:    humidity,
		WindSpeed:   wind,
	}, nil
}

// NewMCPServer 获取MCP服务实例方法
func NewMCPServer() *server.MCPServer {
	// 1. 初始化天气API客户端
	weatherClient := NewWeatherAPIClient()
	// 2. 注册MCP工具（只注册查询天气的MCP）
	mcpServer := server.NewMCPServer(
		"weather-query-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)
	mcpServer.AddTool(
		mcp.NewTool(
			"get_weather",
			mcp.WithDescription("获取指定城市的天气信息"),
			mcp.WithString(
				"city",
				mcp.Description("城市名称，如Beijing、上海"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := request.GetArguments()
			city, ok := args["city"].(string)
			if !ok || city == "" {
				return nil, fmt.Errorf("invalid city parameter")
			}
			weather, err := weatherClient.GetWeather(ctx, city)
			if err != nil {
				return nil, err
			}
			resultText := fmt.Sprintf("城市: %s\n温度: %.1f°C\n天气; %s\n湿度: %d%%\n风速: %.1fkm/h",
				weather.Location,
				weather.Temperature,
				weather.Condition,
				weather.Humidity,
				weather.WindSpeed,
			)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: resultText,
					},
				},
			}, nil
		},
	)
	return mcpServer
}
