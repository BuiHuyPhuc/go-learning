package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Database []struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`
	Security struct {
		Jwt struct {
			Key string `mapstructure:"key"`
		} `mapstructure:"jwt"`
	} `mapstructure:"security"`
}

func main() {
	viper := viper.New()
	viper.AddConfigPath("./config/")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Failed to read configuration %w\n", err))
	}

	fmt.Println("Server port:", viper.GetInt("server.port"))
	fmt.Println("Security jwt:", viper.GetString("security.jwt.key"))

	// config structure
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("Failed to unmarshal configuration %w\n", err))
	}
	fmt.Println("Config port:", config.Server.Port)
	fmt.Println("Config security jwt:", config.Security.Jwt.Key)

	for _, db := range config.Database {
		fmt.Printf("Config database user: %s, pass: %s, host: %s, name: %s\n", db.User, db.Password, db.Host, db.Name)
	}
}
