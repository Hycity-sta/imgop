package routers

import (
	"github.com/gin-gonic/gin"
)

func noRouterHandle(c *gin.Context) {
	c.JSON(404, gin.H{"message": "404 not found"})
}

func Setup() *gin.Engine {
	app := gin.Default()
	SetupApiRouters(app)
	return app
}
