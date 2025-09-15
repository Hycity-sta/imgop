package routers

import (
	"imgop/services"

	"github.com/gin-gonic/gin"
)

func noRouterHandle(c *gin.Context) {
	c.JSON(404, gin.H{"message": "404 not found"})
}

func Setup() *gin.Engine {
	r := gin.Default()

	r.GET("/index", services.GetIndex)

	r.NoRoute(noRouterHandle)

	return r
}
