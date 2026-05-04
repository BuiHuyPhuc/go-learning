package initialize

import (
	"go-learning/global"
	"go-learning/internal/database"
	"go-learning/internal/service"
	"go-learning/internal/service/impl"
)

func InitServiceInterface() {
	queries := database.New(global.Mdbc)

	// User Service Interface
	service.InitUserLogin(impl.NewUserLoginImpl(queries))
}
