package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"imgop/internal/middlewares"
	"imgop/internal/services"
)

func SetupApiRouters(r *gin.Engine) {
	r.NoRoute(noRouterHandle)

	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	public := r.Group("/api")
	public.POST("/login", services.Login)
	public.POST("/signup", services.Signup)
	public.POST("/jwt-auth", services.JwtAuth)

	private := r.Group("/api")
	private.Use(middlewares.JWTAuthMiddleware())

	private.POST("/add-friend", services.AddFriend)
	private.GET("/outgoing-friendrequests", services.GetOutgoingFriendRequests)
	private.GET("/incoming-friendrequests", services.GetIncomingFriendRequests)
	private.POST("/accept-friendrequests", services.AcceptFriendRequests)
	private.POST("/reject-friendrequests", services.RejectFriendRequests)
	private.GET("/friends", services.GetAllFriends)
	private.POST("/search-friend", services.SearchFriend)

}
