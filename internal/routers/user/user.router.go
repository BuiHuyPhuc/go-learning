package user

import "github.com/gin-gonic/gin"

type UserRouter struct {
}

func (ur *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// public router
	userRouterPublic := Router.Group("/users")
	{
		userRouterPublic.POST("/register") // register -> YES -> No
		userRouterPublic.POST("/opt")
	}

	// private router
	userRouterPrivate := Router.Group("/users")
	// userRouterPrivate.Use(LimiterMiddleware(), AuthMiddleware(), PermissionMiddleware())
	{
		userRouterPrivate.GET("/get-info")
	}
}
