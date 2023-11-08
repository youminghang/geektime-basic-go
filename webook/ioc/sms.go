package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/memory"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	// 换内存，还是换别的
	//svc := ratelimit.NewRatelimitSMSService(memory.NewService(),
	//	limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100))
	//return retryable.NewService(svc, 3)
	// 接入监控
	//return metrics.NewPrometheusDecorator(memory.NewService())
	return memory.NewService()
}
