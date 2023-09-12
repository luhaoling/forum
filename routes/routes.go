package routes

import (
	"net/http"
	"project/controllers"
	"project/logger"
	"project/middleware"

	_ "project/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUp(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/api/v1")

	// 注册
	v1.POST("/signup", controllers.SignUpHandler)

	// 登录
	v1.POST("/login", controllers.LoginHandler)

	// 返回帖子列表（支持分页）
	v1.GET("/posts", controllers.GetPostListHandler)

	// 根据时间或分数返回帖子列表（可指定社区）（支持分页）
	v1.GET("/posts2", controllers.GetPostListHandler2)

	// 根据帖子 ID 获取指定帖子
	v1.GET("/post/:id", controllers.GetPostDetailHandler)

	// 获取社区列表
	v1.GET("/community", controllers.CommunityHandler)

	// 根据社区 id 获取社区详情
	v1.GET("/community/:id", controllers.CommunityDetailHandler)

	v1.Use(middleware.JWTAuthMiddleware()) // 应用 JWT 认证中间件
	{
		// 创建帖子
		v1.POST("/post", controllers.CreatePostHandler)

		// 投票
		v1.POST("/vote", controllers.PostVoteController)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
