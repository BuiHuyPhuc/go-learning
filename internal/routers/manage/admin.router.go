package manage

import "github.com/gin-gonic/gin"

type AdminRouter struct {
}

func (ar *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {
	// public router
	adminRouterPublic := Router.Group("/admins")
	{
		adminRouterPublic.POST("/login")
	}

	// private router
	adminRouterPrivate := Router.Group("/admins/users")
	// userRouterPrivate.Use(LimiterMiddleware(), AuthMiddleware(), PermissionMiddleware())
	{
		adminRouterPrivate.POST("/active-user")
	}
}
