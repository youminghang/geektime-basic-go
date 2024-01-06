package grpc

import (
	"context"
	"github.com/go-kratos/aegis/circuitbreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreakerInterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
}

func (b *CircuitBreakerInterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if b.breaker.Allow() == nil {
			res, err := handler(ctx, req)
			if err == nil {
				b.breaker.MarkSuccess()
			} else {
				b.breaker.MarkFailed()
			}
			return res, err
		}
		b.breaker.MarkFailed()
		// 这边你可以考虑使用 Unavailable，或者自己定义一个错误码
		return nil, status.Errorf(codes.Unavailable, "触发了熔断")
	}
}
