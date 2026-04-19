package initialize

import (
	"fmt"
	"go-learning/global"

	"github.com/spf13/viper"
)

func LoadConfig() {
	// TODO: Load config from file or env
	viper := viper.New()
	viper.AddConfigPath("./config/")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Failed to read configuration %w\n", err))
	}

	// fmt.Println("Server port:", viper.GetInt("server.port"))
	// fmt.Println("Security jwt:", viper.GetString("security.jwt.key"))

	// config structure
	if err := viper.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("Unable to decode configuration %w\n", err))
	}
}
