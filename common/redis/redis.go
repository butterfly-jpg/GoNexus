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

// DeleteRedisIndex 删除 Redis 索引
func DeleteRedisIndex(ctx context.Context, filename string) error {
	// 1. 生成索引
	indexName := generateIndexName(filename)
	// 2. 删除对应的索引
	if err := Rdb.Do(ctx, "FT.DROPINDEX", indexName).Err(); err != nil {
		return fmt.Errorf("delete index failed. err: %v", err)
	}
	fmt.Println("delete index success. index name: ", indexName)
	return nil
}

// InitRedisIndex 初始化 Redis 索引
func InitRedisIndex(ctx context.Context, filename string, dimension int) error {
	// 1. 检查索引是否存在,存在就跳过
	indexName := generateIndexName(filename)
	_, err := Rdb.Do(ctx, "FT.INFO", indexName).Result()
	if err == nil {
		fmt.Println("Index already exists, skip creation.")
		return nil
	}
	// 2. 如果索引不存在就创建新索引
	if !strings.Contains(err.Error(), "Unknow index name") {
		return fmt.Errorf("check Index Failure: %w", err)
	}
	fmt.Println("Creating Index:", indexName)
	// 以文件名的前缀作为索引
	prefix := GenerateIndexNamePrefix(filename)
	createArgs := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH",
		"PREFIX", "1", prefix,
		"SCHEMA",
		"content", "TEXT",
		"metadata", "TEXT",
		"vector", "VECTOR", "FLAT",
		"6",
		"TYPE", "FLOAT32",
		"DIM", dimension,
		"DISTANCE_METRIC", "COSINE",
	}
	if err = Rdb.Do(ctx, createArgs...).Err(); err != nil {
		return fmt.Errorf("create index failed. err: %v", err)
	}
	fmt.Println("create index success.")
	return nil
}
