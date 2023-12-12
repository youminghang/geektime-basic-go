package main

import (
	"gitee.com/geekbang/basic-go/webook/pkg/ginx"
	"gitee.com/geekbang/basic-go/webook/pkg/grpcx"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperV2Watch()
	initPrometheus()
	app := Init()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	go func() {
		err := app.migratorServer.Start()
		panic(err)
	}()
	err := app.server.Serve()
	panic(err)
}

func initViperV2Watch() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

type App struct {
	server         *grpcx.Server
	migratorServer *ginx.Server
	consumers      []saramax.Consumer
}
