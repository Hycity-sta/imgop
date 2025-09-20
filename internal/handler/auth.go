package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"imgop/internal/models"
	"imgop/internal/services"
)

func Signup(c *gin.Context) {
	var user models.User

	// 绑定 JSON 到结构体
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateUser(&user)

	// 处理邮箱重复错误
	if mongo.IsDuplicateKeyError(err) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "此邮箱已经注册",
			"code":  "EMAIL_EXISTS",
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"name":          user.Name,
			"password_hash": user.PasswordHash,
			"email":         user.Email,
		})
	}

}

func Login(c *gin.Context) {

}
