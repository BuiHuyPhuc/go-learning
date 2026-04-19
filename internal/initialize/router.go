package initialize

import (
	"go-learning/global"
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
	manageRouter := routers.RouterGroupApp.Manage
	userRouter := routers.RouterGroupApp.User

	MainGroup := r.Group("/api/v1")
	{
		MainGroup.GET("/check-status") // tracking monitor
	}
	{
		manageRouter.InitAdminRouter(MainGroup)
		manageRouter.InitUserRouter(MainGroup)
	}
	{
		userRouter.InitUserRouter(MainGroup)
		userRouter.InitProductRouter(MainGroup)
	}

	return r
}
