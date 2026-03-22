package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"GoNexus/config"

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
