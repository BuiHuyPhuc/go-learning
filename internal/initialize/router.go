package initialize

import (
	"go-learning/global"
	"go-learning/internal/middlewares"
	"go-learning/internal/routers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// TODO: Init router from file or env
	var r *gin.Engine

	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New() // Không ghi nhật ký
	}

	// middlewares
	// r.Use() // logging middleware
	// r.Use() // cross-domain middleware
	// r.Use() // limiter middleware
	r.Use(middlewares.ValidationMiddleware()) // validation middleware

	// vegeta attack load testing tool
	// r.Use(middlewares.NewRateLimiter().GlobalRateLimiter())
	// r.GET("/ping/100", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "pong 100",
	// 	})
	// })

	// r.Use(middlewares.NewRateLimiter().PublicRateLimiter())
	// r.GET("/ping/80", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "pong 80",
	// 	})
	// })

	// r.Use(middlewares.NewRateLimiter().UserPrivateRateLimiter())
	// r.GET("/ping/50", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "pong 50",
	// 	})
	// })

	// manageRouter := routers.RouterGroupApp.Manage
	userRouter := routers.RouterGroupApp.User

	MainGroup := r.Group("/api/v1")
	{
		MainGroup.GET("/check-status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "healthy",
			})
		}) // tracking monitor
	}
	{
		// manageRouter.InitAdminRouter(MainGroup)
		// manageRouter.InitUserRouter(MainGroup)
	}
	{
		userRouter.InitUserRouter(MainGroup)
		// userRouter.InitProductRouter(MainGroup)
		userRouter.InitTicketRouter(MainGroup)
	}

	return r
}
