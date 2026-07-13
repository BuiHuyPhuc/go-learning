package initialize

import (
	"go-learning/global"
	"go-learning/internal/database"
	"go-learning/internal/service"
	"go-learning/internal/service/impl"
)

func InitServiceInterface() {
	queries := database.New(global.Mdbc)

	if global.Rdb == nil {
		panic("global.Rdb is nil! Redis not initialized")
	}

	redisCache := impl.NewRedisCacheImpl(global.Rdb)

	// User Service Interface
	service.InitUserLogin(impl.NewUserLoginImpl(queries))
	// Ticket Service Interface
	service.InitTicketItem(impl.NewTicketItemImpl(queries, redisCache))
}
