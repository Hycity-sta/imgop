package services

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"imgop/internal/models"
	"imgop/internal/utils"
)

type JwtAuthRequest struct {
	Token string `json:"token" binding:"required"`
}

func JwtAuth(c *gin.Context) {
	var req JwtAuthRequest

	// 绑定 JSON 到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务端接收请求时出错"})
		return
	}

	// jwttoken解码
	res, err := utils.ParseToken(req.Token)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token 过期"})
		case errors.Is(err, jwt.ErrTokenMalformed):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token 格式错误"})
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token 验证失败"})
		}
	}

	// 数据库查询token指定的用户存不存在
	id := res.UserID
	user := models.FindUserById(id)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token失效, 用户不存在"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "ok"})
	}
}
