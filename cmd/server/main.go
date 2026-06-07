package main

import (
	_ "go-learning/cmd/swag/docs"
	"go-learning/internal/initialize"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var pingCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "ping_request_count_total",
		Help: "Total number of ping requests.",
	},
)

func ping(c *gin.Context) {
	pingCounter.Inc() // 1 2 3 ...
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

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

	prometheus.MustRegister(pingCounter)

	r.GET("/ping/200", ping)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8888"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
