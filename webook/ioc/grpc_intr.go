package ioc

import (
	intrv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/intr/v1"
	"gitee.com/geekbang/basic-go/webook/interactive/service"
	"gitee.com/geekbang/basic-go/webook/internal/web/client"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitIntrGRPCClientV1 这个版本是使用 ETCD 做服务发现的版本
func InitIntrGRPCClientV1(l logger.LoggerV1) intrv1.InteractiveServiceClient {
	type Config struct {
		Name     string `yaml:"name"`
		EtcdAddr string `yaml:"etcdAddr"`
		Secure   bool   `yaml:"secure"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	cli, err := clientv3.NewFromURL("http://localhost:12379")
	if err != nil {
		panic(err)
	}
	etcdResolver, err := resolver.NewBuilder(cli)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(etcdResolver)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial("etcd:///service/"+cfg.Name, opts...)
	if err != nil {
		panic(err)
	}
	return intrv1.NewInteractiveServiceClient(cc)
}

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
	local := client.NewInteractiveServiceAdapter(svc)
	res := client.NewInteractiveClient(remote, local, cfg.Threshold)

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
