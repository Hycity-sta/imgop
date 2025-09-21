package routers

import (
	"github.com/gin-gonic/gin"
)

func noRouterHandle(c *gin.Context) {
	c.JSON(404, gin.H{"message": "404 not found"})
}

func Setup() *gin.Engine {
	r := gin.Default()

	r.NoRoute(noRouterHandle)

	auth_ := r.Group("/auth")
	{
		auth_.POST("/login")
		auth_.POST("/signup")
	}

	user_ := r.Group("/user")
	{
		user_.GET("/GetAllFriends")
		user_.GET("/GetAllGroups")
		user_.POST("/AddNewFriendById")
		user_.POST("/JoinNewGroupById")
		user_.GET("/FindFriendByName")
		user_.GET("/FindGroupByName")
		user_.POST("/DeleteFriendById")
		user_.POST("/DeleteGroupById")
	}

	chat_ := r.Group("/chat")
	{
		chat_.GET("/")
	}

	return r
}
