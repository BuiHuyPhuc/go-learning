package initialize

import (
	"go-learning/global"
	"go-learning/pkg/logger"
)

func InitLogger() {
	// TODO: Init logger from file or env
	global.Logger = logger.NewLogger(global.Config.Logger)
}
