package global

import (
	"go-learning/pkg/logger"
	"go-learning/pkg/settings"

	"gorm.io/gorm"
)

var (
	Config settings.Config
	Logger *logger.LoggerZap
	Mdb    *gorm.DB
)

/*
Config
Mysql
Redis
*/
