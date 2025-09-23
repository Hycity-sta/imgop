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
	private.GET("/friends", services.GetAllFriends)
	private.POST("/search-friend", services.SearchFriend)
	private.POST("/add-friend", services.AddFriend)
	private.POST("/accept-friend", services.AcceptFriend)
	private.POST("/reject-friend", services.RejectFriend)
	private.GET("/pending-friends", services.GetAllPendingFriend)

}
