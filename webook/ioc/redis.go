package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	// 这里演示读取特定的某个字段
	cmd := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
	})
	return cmd
}

func InitRLockClient(cmd redis.Cmdable) *rlock.Client {
	return rlock.NewClient(cmd)
}
