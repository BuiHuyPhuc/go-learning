package initialize

import (
	"go-learning/global"
	"go-learning/pkg/logger"

	"go.uber.org/zap"
)

func InitLogger() {
	// TODO: Init logger from file or env
	global.Logger = logger.NewLogger(global.Config.Logger)
	global.Logger.Info("Logger initialized", zap.String("ok", "success"))
}
