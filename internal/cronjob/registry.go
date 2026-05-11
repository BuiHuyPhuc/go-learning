package cronjob

import "go-learning/global"

func RegistryRunCron() {
	SendEmailForVIPEvery3Seconds()
	GetInfoUserEvery5Seconds()

	global.Cron.Start()
}
