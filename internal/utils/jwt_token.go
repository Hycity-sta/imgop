package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义 Claims 结构体，包含标准 Claims 和自定义字段
type CustomClaims struct {
	jwt.RegisteredClaims
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

var (
	SecretKey = []byte("123") // 密钥
	Issuer    = "imgop"       // 签发者名称
)

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string) (string, error) {
	// 设置 Claims
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                     // 生效时间
			Issuer:    Issuer,                                             // 签发者
		},
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整 Token
	return token.SignedString(SecretKey)
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证 Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
