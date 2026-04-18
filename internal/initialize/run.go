package initialize

import (
	"fmt"
	"go-learning/global"
	"log"

	"go.uber.org/zap"
)

func Run() {
	// load configuration
	LoadConfig()
	fmt.Printf("Loading configuration mysql:\n%+v\n", global.Config.Mysql)

	InitLogger()
	global.Logger.Info("Logger initialized", zap.String("ok", "success"))

	InitMysql()
	InitRedis()

	r := InitRouter()

	// Server will listen on 0.0.0.0:8888 (localhost:8888 on Windows)
	if err := r.Run(":8888"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
