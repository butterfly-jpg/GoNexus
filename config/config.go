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

// EmailConfig 邮箱的配置结构体
type EmailConfig struct {
	Authcode string `toml:"authcode"`
	Email    string `toml:"email"`
}

type RedisKeyConfig struct {
	CaptchaPrefix string
}

var DefaultRedisKeyConfig = RedisKeyConfig{
	CaptchaPrefix: "captcha:%s",
}

// Config 配置结构体
type Config struct {
	MysqlConfig `toml:"mysqlConfig"`
	RedisConfig `toml:"redisConfig"`
	EmailConfig `toml:"emailConfig"`
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
