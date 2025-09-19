package routers

import (
	"github.com/gin-gonic/gin"

	"imgop/handler"
)

func noRouterHandle(c *gin.Context) {
	c.JSON(404, gin.H{"message": "404 not found"})
}

func Setup() *gin.Engine {
	r := gin.Default()

	r.GET("/index", handler.GetIndex)
	r.NoRoute(noRouterHandle)

	return r
}
