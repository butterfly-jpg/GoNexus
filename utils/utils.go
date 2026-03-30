package utils

import (
	"GoNexus/model"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"GoNexus/config"

	"github.com/cloudwego/eino/schema"
	"github.com/golang-jwt/jwt/v5"
)

// GetRandomNumbers 获取随机位数的数字
func GetRandomNumbers(num int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := ""
	for i := 0; i < num; i++ {
		digit := r.Intn(10)
		code += strconv.Itoa(digit)
	}
	return code
}

// MD5 MD5加密
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// Claims JWT的声明结构体
type Claims struct {
	UserID   int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 根据用户ID和用户名生成JWT token
func GenerateToken(userID int64, username string) (string, error) {
	cfg := config.GetConfig()
	secret := cfg.JWTConfig.Secret
	if secret == "" {
		secret = "gonexus_default_secret"
	}
	expireDay := cfg.JWTConfig.ExpireDay
	if expireDay <= 0 {
		expireDay = 7
	}
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, expireDay)),
			Issuer:    cfg.JWTConfig.Issuer,
			Subject:   cfg.JWTConfig.Subject,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析JWT token
func ParseToken(tokenStr string) (string, error) {
	cfg := config.GetConfig()
	secret := cfg.JWTConfig.Secret
	if secret == "" {
		secret = "gonexus_default_secret"
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Username, nil
}

// ConvertToSchemaMessages 将model.Message转换为schema.Message格式
func ConvertToSchemaMessages(msgs []*model.Message) []*schema.Message {
	schemaMsgs := make([]*schema.Message, 0, len(msgs))
	for _, msg := range msgs {
		role := schema.Assistant
		if msg.IsUser {
			role = schema.User
		}
		schemaMsgs = append(schemaMsgs, &schema.Message{
			Role:    role,
			Content: msg.Content,
		})
	}
	return schemaMsgs
}

// ConvertToModelMessages 将schema.Message转换为model.Message格式
func ConvertToModelMessages(sessionID, username, content string) *model.Message {
	return &model.Message{
		SessionID: sessionID,
		Username:  username,
		Content:   content,
	}
}

// ValidateFile 校验文件是否为允许的本地文件,只支持.md和.txt文件
func ValidateFile(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".md" && ext != ".txt" {
		return fmt.Errorf("file type must be .md or .txt, current file type is %s", ext)
	}
	return nil
}

// RemoveAllFilesInDir 删除目录中的所有文件（不删除子目录）
func RemoveAllFilesInDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// 说明连目录都不存在也就不需要再删除文件了
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if err = os.Remove(filepath.Join(dir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}
