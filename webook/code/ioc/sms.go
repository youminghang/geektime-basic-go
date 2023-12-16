package ioc

import (
	smsv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/sms/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitSmsRpcClient() smsv1.SmsServiceClient {
	type config struct {
		Addr string `yaml:"addr"`
	}
	var cfg config
	err := viper.UnmarshalKey("smsGrpc", &cfg)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(cfg.Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := smsv1.NewSmsServiceClient(conn)
	return client
}
