package main

import (
	"GoNexus/common/aihelper"
	"GoNexus/common/mysql"
	"GoNexus/common/rabbitmq"
	"GoNexus/common/redis"
	"GoNexus/config"
	"GoNexus/dao/message"
	"GoNexus/router"
	"fmt"
	"log"
)

// StartServer 开启服务
func StartServer(addr string, port int) error {
	// 1. 初始化路由
	r := router.InitRouter()
	// 2. 加载服务器静态资源路径映射关系
	return r.Run(fmt.Sprintf("%s:%d", addr, port))
}

// readDataFromDB 初始化消息并初始化 AIHelperManager
func readDataFromDB() error {
	messages, err := message.GetAllMessages()
	if err != nil {
		return err
	}
	globalManager := aihelper.GetGlobalManager()
	for _, msg := range messages {
		helper, err := globalManager.GetOrCreateAIHelper(msg.Username,
			msg.SessionID, "1", make(map[string]interface{}))
		if err != nil {
			log.Printf("[readDataFromDB] failed to create helper for user=%s session=%s: .err: %v",
				msg.Username, msg.SessionID, err)
			continue
		}
		log.Printf("[readDataFromDB] init: %s", msg.SessionID)
		// 添加消息到存储中
		helper.AddMessage(msg.Content, msg.Username, msg.IsUser, false)
	}
	return nil
}

func main() {
	// 1. 获取GoNexus后端服务IP地址和端口号
	conf := config.GetConfig()
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port
	// 2. 初始化MySQL
	if err := mysql.InitMysql(); err != nil {
		log.Println("init MySQL failed. err:", err)
		return
	}
	log.Println("init MySQL success")
	// 3. 初始化 AIHelperManager
	readDataFromDB()
	log.Println("init AIHelperManager success")
	// 3. 初始化Redis
	redis.InitRedis()
	log.Println("init Redis success")
	// 4. 初始化RabbitMQ
	rabbitmq.InitRabbitMQ()
	log.Println("init RabbitMQ success")
	// 5. 注册并启动HTTP服务
	if err := StartServer(host, port); err != nil {
		log.Println("start server failed. err:", err)
		return
	}
}
