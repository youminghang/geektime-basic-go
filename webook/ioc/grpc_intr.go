package ioc

import (
	intrv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/intr/v1"
	client2 "gitee.com/geekbang/basic-go/webook/bff/client"
	"gitee.com/geekbang/basic-go/webook/interactive/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitIntrGRPCClient(svc service.InteractiveService,
	l logger.LoggerV1) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string
		Secure    bool
		Threshold int32
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}

	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client2.NewInteractiveServiceAdapter(svc)
	res := client2.NewInteractiveClient(remote, local, cfg.Threshold)

	viper.OnConfigChange(func(in fsnotify.Event) {
		// 重置整个 Config
		cfg = Config{}
		err1 := viper.UnmarshalKey("grpc.intr", cfg)
		if err1 != nil {
			l.Error("重新加载grpc.intr的配置失败", logger.Error(err1))
			return
		}
		// 这边更新 Threshold
		res.UpdateThreshold(cfg.Threshold)
	})
	return res
}
