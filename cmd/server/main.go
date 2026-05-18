package main

import (
	"go-learning/internal/initialize"
	"log"

	_ "github.com/BuiHuyPhuc/go-learning/cmd/swag/docs"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           API Documentation Ecommerce Backend SHOPDEVGO
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  https://github.com/BuiHuyPhuc/go-learning

// @contact.name   TEAM PHUCBUI
// @contact.url    https://github.com/BuiHuyPhuc/go-learning
// @contact.email  buihuyphuc97@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8888
// @BasePath  /api/v1
// @schemes   http
func main() {
	r := initialize.Run()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8888"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
