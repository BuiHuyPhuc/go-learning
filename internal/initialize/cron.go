package initialize

import (
	"go-learning/global"
	"go-learning/internal/cronjob"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	c := cron.New(cron.WithSeconds())
	global.Logger.Info("Initializing Cron successfully")
	global.Cron = c

	cronjob.RegistryRunCron()
}
