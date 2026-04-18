package initialize

import (
	"fmt"
	c "go-learning/internal/controller"
	"go-learning/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func AA() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Before -> AA")
		c.Next()
		fmt.Println("After -> AA")
	}
}

func BB() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Before -> BB")
		c.Next()
		fmt.Println("After -> BB")
	}
}

func CC(c *gin.Context) {
	fmt.Println("Before -> CC")
	c.Next()
	fmt.Println("After -> CC")
}

func InitRouter() *gin.Engine {
	// TODO: Init router from file or env
	r := gin.Default()

	r.Use(middlewares.AuthMiddleware(), AA(), BB(), CC)

	v1 := r.Group("/api/v1")
	v1.GET("/ping", c.NewPongController().Pong)
	v1.GET("/users/1", c.NewUserController().GetUserById)

	return r
}
