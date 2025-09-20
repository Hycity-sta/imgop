package routers

import (
	"github.com/gin-gonic/gin"

	"imgop/internal/handler"
)

func noRouterHandle(c *gin.Context) {
	c.JSON(404, gin.H{"message": "404 not found"})
}

func Setup() *gin.Engine {
	r := gin.Default()

	r.NoRoute(noRouterHandle)

	auth_ := r.Group("/auth")
	{
		auth_.GET("/login", handler.Login)
		auth_.GET("/signup", handler.Signup)
	}

	user_ := r.Group("/user")
	{
		user_.GET("/GetFriendList")
	}

	chat_ := r.Group("/chat")
	{
		chat_.GET("/")
	}

	return r
}
