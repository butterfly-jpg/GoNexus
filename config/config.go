package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// MysqlConfig mysql的配置结构体
type MysqlConfig struct {
	MysqlPort         int    `toml:"port"`
	MysqlHost         string `toml:"host"`
	MysqlUser         string `toml:"user"`
	MysqlPassword     string `toml:"password"`
	MysqlDatabaseName string `toml:"databaseName"`
	MysqlCharset      string `toml:"charset"`
}

// RedisConfig redis的配置结构体
type RedisConfig struct {
	RedisPort int    `toml:"port"`
	RedisDb   int    `toml:"db"`
	RedisHost string `toml:"host"`
	RedisPass string `toml:"password"`
}

// RabbitmqConfig RabbitMQ的配置结构体
type RabbitmqConfig struct {
	RabbitmqHost     string `toml:"host"`
	RabbitmqPort     int    `toml:"port"`
	RabbitmqUsername string `toml:"username"`
	RabbitmqPassword string `toml:"password"`
	RabbitmqVhost    string `toml:"vhost"`
}

// RagModelConfig Rag模型配置结构体
type RagModelConfig struct {
	RagEmbeddingModel string `toml:"embeddingModel"`
	RagChatModelName  string `toml:"chatModelName"`
	RagBaseUrl        string `toml:"baseUrl"`
	RagDimension      int    `toml:"dimension"`
}

// EmailConfig 邮箱的配置结构体
type EmailConfig struct {
	Authcode string `toml:"authcode"`
	Email    string `toml:"email"`
}

type RedisKeyConfig struct {
	CaptchaPrefix   string
	IndexName       string
	IndexNamePrefix string
}

var DefaultRedisKeyConfig = RedisKeyConfig{
	CaptchaPrefix:   "captcha:%s",
	IndexName:       "rag_docs:%s:idx",
	IndexNamePrefix: "rag_docs:%s:",
}

// JWTConfig JWT的配置结构体
type JWTConfig struct {
	Secret    string `toml:"secret"`
	ExpireDay int    `toml:"expireDay"`
	Issuer    string `toml:"issuer"`
	Subject   string `toml:"subject"`
}

// MainConfig 后端服务配置信息
type MainConfig struct {
	AppName string `toml:"appName"`
	Port    int    `toml:"port"`
	Host    string `toml:"host"`
}

// Config 配置结构体
type Config struct {
	MainConfig     `toml:"mainConfig"`
	MysqlConfig    `toml:"mysqlConfig"`
	RedisConfig    `toml:"redisConfig"`
	RabbitmqConfig `toml:"rabbitmqConfig"`
	EmailConfig    `toml:"emailConfig"`
	JWTConfig      `toml:"jwtConfig"`
	RagModelConfig `toml:"ragModelConfig"`
}

var config *Config

// InitConfig 初始化配置
func InitConfig() error {
	if _, err := toml.DecodeFile("config/config.toml", config); err != nil {
		log.Fatalf("decode config file failed, err:%v", err)
		return err
	}
	return nil
}

// GetConfig 获取配置信息
func GetConfig() *Config {
	if config == nil {
		config = &Config{}
		_ = InitConfig()
	}
	return config
}
