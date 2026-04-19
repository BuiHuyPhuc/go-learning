package manage

import "github.com/gin-gonic/gin"

type UserRouter struct {
}

func (ur *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// public router

	// private router
	userRouterPrivate := Router.Group("/admins/users")
	// userRouterPrivate.Use(LimiterMiddleware(), AuthMiddleware(), PermissionMiddleware())
	{
		userRouterPrivate.POST("/active-user")
	}
}
