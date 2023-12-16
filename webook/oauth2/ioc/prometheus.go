package ioc

import (
	wechat2 "gitee.com/geekbang/basic-go/webook/internal/oauth2/service/wechat"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/spf13/viper"
)

func InitPrometheus(logv1 logger.LoggerV1) wechat2.Service {
	svc := InitService(logv1)
	type Config struct {
		NameSpace  string `yaml:"nameSpace"`
		Subsystem  string `yaml:"subsystem"`
		InstanceID string `yaml:"instanceId"`
		Name       string `yaml:"name"`
	}
	var cfg Config
	err := viper.UnmarshalKey("prometheus", &cfg)
	if err != nil {
		panic(err)
	}
	return wechat2.NewPrometheusDecorator(svc, cfg.NameSpace, cfg.Subsystem, cfg.InstanceID, cfg.Name)
}
