package ioc

import (
	grpc2 "gitee.com/geekbang/basic-go/webook/interactive/grpc"
	"gitee.com/geekbang/basic-go/webook/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(intr *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	intr.Register(server)
	return &grpcx.Server{
		Server: server,
		Addr:   cfg.Addr,
	}
}
