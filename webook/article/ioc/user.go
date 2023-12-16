package ioc

import (
	userv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/user/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitUserRpcClient() userv1.UserServiceClient {
	type config struct {
		Addr string `yaml:"addr"`
	}
	var cfg config
	err := viper.UnmarshalKey("userGrpc", &cfg)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(cfg.Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := userv1.NewUserServiceClient(conn)
	return client
}
