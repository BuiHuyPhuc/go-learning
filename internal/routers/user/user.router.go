package user

import (
	"go-learning/internal/controller/account"
	"go-learning/internal/wire"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (urt *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// This is non-dependency
	// ur := repo.NewUserRepository()
	// us := service.NewUserService(ur)
	// userHandlerNonDependency := controller.NewUserController(us)

	// Wire go
	userController, _ := wire.InitUserRouterHandler()

	// public router
	userRouterPublic := Router.Group("/users")
	{
		userRouterPublic.POST("/register", userController.Register) // register -> YES -> No
		userRouterPublic.POST("/login", account.Login.Login)
	}

	// private router
	userRouterPrivate := Router.Group("/users")
	// userRouterPrivate.Use(LimiterMiddleware(), AuthMiddleware(), PermissionMiddleware())
	{
		userRouterPrivate.GET("/get-info")
	}
}
