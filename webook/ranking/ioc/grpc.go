package ioc

import (
	"gitee.com/geekbang/basic-go/webook/pkg/grpcx"
	grpc2 "gitee.com/geekbang/basic-go/webook/ranking/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(rankingServer *grpc2.RankingServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	rankingServer.Register(server)
	return &grpcx.Server{
		Server: server,
		Addr:   cfg.Addr,
	}
}
