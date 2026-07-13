package initialize

import (
	"fmt"
	"go-learning/global"

	"github.com/gin-gonic/gin"
)

func Run() *gin.Engine {
	// load configuration
	LoadConfig()
	fmt.Printf("Loading configuration mysql:\n%+v\n", global.Config.Mysql)

	InitLogger()
	// InitMysql()
	InitMysqlC() // docker exec -it mysql bash
	InitRedis()  // docker exec -it redis redis-cli
	// InitRedisSentinel()
	// InitKafka()

	// InitCron()

	InitServiceInterface()

	return InitRouter()

	// Server will listen on 0.0.0.0:8888 (localhost:8888 on Windows)
	// if err := r.Run(":8888"); err != nil {
	// 	log.Fatalf("failed to run server: %v", err)
	// }
}
