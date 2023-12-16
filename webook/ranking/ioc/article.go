package ioc

import (
	articlev1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/article/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitArticleRpcClient() articlev1.ArticleServiceClient {
	type config struct {
		Addr string `yaml:"addr"`
	}
	var cfg config
	err := viper.UnmarshalKey("articleGrpc", &cfg)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(cfg.Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := articlev1.NewArticleServiceClient(conn)
	return client
}
