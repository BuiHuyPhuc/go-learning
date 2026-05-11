package cronjob

import (
	"fmt"
	"go-learning/global"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func SendEmailForVIPEvery3Seconds() {
	fmt.Println("...SendEmailForVIPEvery3Seconds")

	_, err := global.Cron.AddFunc("*/3 * * * * *", func() {
		fmt.Println("...run... 3 seconds")
	})

	if err != nil {
		global.Logger.Error("cron not active", zap.Error(err))
	}
}

func GetInfoUserEvery5Seconds() {
	fmt.Println("...SendEmailForVIPEvery3Seconds")

	_, err := global.Cron.AddFunc("*/5 * * * * *", func() {
		fmt.Println("...run... 5 seconds")
		res, err := http.Get("https://api.github.com/users/go-learning")
		if err != nil {
			global.Logger.Error("go-learning not active", zap.Error(err))
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			global.Logger.Error("body not active", zap.Error(err))
			return
		}

		fmt.Printf("Github is res: %s", string(body))
	})

	if err != nil {
		global.Logger.Error("cron not active", zap.Error(err))
	}
}
