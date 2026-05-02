package global

import (
	"database/sql"
	"go-learning/pkg/logger"
	"go-learning/pkg/settings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config settings.Config
	Logger *logger.LoggerZap
	Mdb    *gorm.DB
	Mdbc   *sql.DB
	Rdb    *redis.Client
)

/*
Config
Mysql
Redis
*/
