// Package redis 提供 redis 生成key的方法
package redis

import (
	"GoNexus/config"
	"fmt"
)

// generateCaptcha 根据邮箱生成对应的key
func generateCaptcha(email string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.CaptchaPrefix, email)
}

// generateIndexName 根据文件名生成对应的key
func generateIndexName(filename string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.IndexName, filename)
}

// generateIndexNamePrefix 根据文件名生成对应的前缀key
func generateIndexNamePrefix(filename string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.IndexNamePrefix, filename)
}
