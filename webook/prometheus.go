package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		// 监听 8081 端口，你也可以做成可配置的
		http.ListenAndServe(":8081", nil)
	}()
}
