package global

import (
	"database/sql"
	"go-learning/pkg/logger"
	"go-learning/pkg/settings"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

var (
	Config        settings.Config
	Logger        *logger.LoggerZap
	Mdb           *gorm.DB
	Mdbc          *sql.DB
	Rdb           *redis.Client
	KafkaProducer *kafka.Writer
	Cron          *cron.Cron
)

/*
Config
Mysql
Redis
*/
