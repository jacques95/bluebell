package redis

import (
	"bluebell/settings"
	"fmt"
	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Port, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Dd,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	return
}
func Close() {
	_ = rdb.Close()
}
