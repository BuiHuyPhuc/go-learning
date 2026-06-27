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

/*
docker network inspect redis-master_default
docker exec -it redis redis-cli
keys *
getrange, mset, mget, incr, incrby, decr, decrby, expire, ttl
lpush, lrange, rpush, llen, lpop, rpop, lset, linsert, lindex, lpushx (nếu key không tồn tại thì k add vào), sort <key> asc|desc alpha
sadd, smembers, scard, sismember, sdiff, sdiffstore, sinter, sinterstore, sunion, sunionstore
INFO replication


1. Redis transaction: Sử dụng từ khóa multi để bắt đầu trans và kết thúc exec
VD: multi -> hset key field1 value + hget key field2 -> exec
a. Nếu trong quá trình thực hiện trans bị lỗi thì khi exec sẽ "EXECABORT Transaction discarded because of previous errors."
b. Nhưng trong quá trình này nếu bị can thiệp hset key field2 bởi luồng khác thì sau khi exec sẽ hiển thị dữ liệu như đã hset
-> Nếu không muốn bị thay đổi trong trans thì sử dụng "watch key", nếu phát hiện key bị thay đổi thì exec = nil
VD: watch key -> multi -> hset key field1 value + hget key field2 -> exec


2. Sử dụng lệnh LUA
- Không rollback
- Nếu trong quá trình exec bị lỗi thì sẽ trả về exception với redis.call() và bypass với redis.pcall()
- Nếu có nhiều tập lệnh LUA trên cluster hay master sẽ tải vào bộ đệm theo thứ tự hàng đợi không bị gián đoạn bởi client khác nhau nên đảm bảo tính nguyên tử (atomic)

VD:
eval "return ARGV[1]" 0 anonystick
eval "return {ARGV[1] ARGV[2]}" 0 anonystick tipsjava

eval "redis.call('SET', KEYS[1], ARGV[1])" 1 ticket:1 10000
eval "if redis.call('EXISTS', KEYS[1]) == 0 then return redis.call('SET', KEYS[1], ARGV[1]) else return -1 end" 1 ticket:2 50000

eval "redis.call('SET', 'k1', 'v1'); redis.call('INCRBY', 'k2', 1/0); redis.call('SET', 'k3', 'v3')" 0
eval "redis.pcall('SET', 'k4', 'v4'); redis.pcall('INCRBY', 'k5', 1/0); redis.pcall('SET', 'k6', 'v6')" 0
*/
