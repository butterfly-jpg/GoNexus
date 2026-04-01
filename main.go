package main

import (
	"GoNexus/common/mysql"
	"GoNexus/common/rabbitmq"
	"GoNexus/common/redis"
	"GoNexus/config"
	"GoNexus/router"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// StartServer 开启服务
func StartServer(addr string, port int) error {
	// 1. 初始化路由
	r := router.InitRouter()
	// 2. 加载服务器静态资源路径映射关系
	return r.Run(fmt.Sprintf("%s:%d", addr, port))
}

func main() {
	// 0. 加载 .env 环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, using system environment variables")
	}
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
