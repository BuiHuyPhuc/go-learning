package initialize

import (
	"context"
	"fmt"
	"go-learning/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func InitRedis() {
	// TODO: Init redis from file or env
	r := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", r.Host, r.Port),
		Password: r.Password, // no password set
		DB:       r.DB,       // use default DB
		PoolSize: 10,         // use default PoolSize
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Error("Redis initialization error", zap.Error(err))
	}
	global.Logger.Info("Initializing Redis successfully")
	global.Rdb = rdb
	// redisExample()
}

func redisExample() {
	err := global.Rdb.Set(ctx, "score", 100, 0).Err()
	if err != nil {
		fmt.Println("Error redis set", zap.Error(err))
		return
	}

	val, err := global.Rdb.Get(ctx, "score").Result()
	if err != nil {
		fmt.Println("Error redis get", zap.Error(err))
		return
	}

	global.Logger.Info("Value score is::", zap.String("score", val))
}
