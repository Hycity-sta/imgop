package services

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"imgop/internal/models"
	"imgop/internal/utils"
)

// 接收前端请求的结构体
type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 注册服务
func Signup(c *gin.Context) {
	var req SignupRequest

	// 绑定 JSON 到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "绑定结构体出错"})
		return
	}

	// 密码hash加密
	hash, err1 := utils.HashPassword(req.Password)
	if err1 != nil {
		log.Println(err1)
	}

	// 创建user对象
	user := models.User{
		Name:         req.Name,
		PasswordHash: hash,
		Email:        req.Email,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 数据库插入新的user
	err := models.InsertUser(user)

	// 处理邮箱重复错误
	if mongo.IsDuplicateKeyError(err) {
		c.JSON(http.StatusConflict, gin.H{"error": "此邮箱已经注册"})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"name":          user.Name,
			"password_hash": user.PasswordHash,
			"email":         user.Email,
		})
	}
}

// 用来接收前端请求的结构体
type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录服务
func Login(c *gin.Context) {
	var req LoginRequest

	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	// 用户验证
	user := models.FindUserByEmail(req.Email)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或邮箱出错误"})
		return
	}

	// 生成 JWT Token
	userID := user.ID.Hex()
	user_email := user.Email
	token, err := utils.GenerateToken(userID, user_email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失败"})
		return
	}

	// 返回 Token
	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"success":  "ok",
		"username": user.Name,
	})
}