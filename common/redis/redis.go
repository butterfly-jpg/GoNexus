package redis

import (
	"GoNexus/config"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis配置
func InitRedis() {
	conf := config.GetConfig()
	port := conf.RedisConfig.RedisPort
	db := conf.RedisConfig.RedisDb
	host := conf.RedisConfig.RedisHost
	password := conf.RedisConfig.RedisPass
	addr := fmt.Sprintf("%s:%s", host, strconv.Itoa(port))

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Protocol: 2,
		Password: password,
		DB:       db,
	})
}

// CheckCaptchaForEmail 检验邮箱对应的验证码缓存是否有效
func CheckCaptchaForEmail(email string, captcha string) (bool, error) {
	key := generateCaptcha(email)
	storedCaptcha, err := Rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	if strings.EqualFold(storedCaptcha, captcha) {
		// 验证成功缓存存在，然后把缓存删掉
		Rdb.Del(ctx, key)
		return true, nil
	}
	return false, nil
}

// SetCaptchaForEmail 为邮箱为key验证码为value设置缓存
func SetCaptchaForEmail(email, captcha string) error {
	key := generateCaptcha(email)
	// 缓存过期时间为两分钟
	expire := 2 * time.Minute
	return Rdb.Set(ctx, key, captcha, expire).Err()
}

// generateCaptcha 根据邮箱生成对应的key
func generateCaptcha(email string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.CaptchaPrefix, email)
}
