package routers

import (
	c "go-learning/internal/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	v1.GET("/ping", c.NewPongController().Pong)
	v1.GET("/users/1", c.NewUserController().GetUserById)

	return r
}
