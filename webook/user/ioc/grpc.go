package ioc

import (
	"gitee.com/geekbang/basic-go/webook/pkg/grpcx"
	grpc2 "gitee.com/geekbang/basic-go/webook/user/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(userServer *grpc2.UserServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	userServer.Register(server)
	return &grpcx.Server{
		Server: server,
		Addr:   cfg.Addr,
	}
}
