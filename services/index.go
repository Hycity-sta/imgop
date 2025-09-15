package services

import (
	"github.com/gin-gonic/gin"
)

func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Your MVC Gin App is running!"})
}