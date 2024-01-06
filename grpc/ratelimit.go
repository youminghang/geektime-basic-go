package grpc

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/pkg/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type InterceptorBuilder struct {
	limiter ratelimit.Limiter
	key     string
}

func (b *InterceptorBuilder) BuildServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		ok, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			// 这里采用保守措施，在触发限流之后直接返回
			return nil, err
		}
		if !ok {
			return nil, status.Errorf(codes.ResourceExhausted,
				"限流")
		}
		return handler(ctx, req)
	}
}

// BuildServerUnaryInterceptorService 服务级别限流
func (b *InterceptorBuilder) BuildServerUnaryInterceptorService() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		// prefix 这里可以做成参数
		if strings.HasPrefix(info.FullMethod, "/UserService") {
			ok, err := b.limiter.Limit(ctx, b.key)
			if err != nil {
				// 这里采用保守措施，在触发限流之后直接返回
				return nil, err
			}
			if !ok {
				return nil, status.Errorf(codes.ResourceExhausted,
					"限流")
			}
		}

		return handler(ctx, req)
	}
}

// BuildServerUnaryInterceptorV1 配合后面的降级逻辑处理
func (b *InterceptorBuilder) BuildServerUnaryInterceptorV1() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		ok, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			// 这里采用保守措施，在触发限流之后直接返回
			return nil, err
		}
		if !ok {
			ctx = context.WithValue(ctx, "downgrade", "true")
			return handler(ctx, req)
		}
		return handler(ctx, req)
	}
}

// BuildServerUnaryInterceptorBiz 针对业务限流
func (b *InterceptorBuilder) BuildServerUnaryInterceptorBiz() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		// prefix 这里可以做成参数
		if getById, ok := req.(*GetByIdReq); ok {
			key := fmt.Sprintf("limiter:user:get_by_id:%d", getById.GetId())
			ok, err := b.limiter.Limit(ctx, key)
			if err != nil {
				// 这里采用保守措施，在触发限流之后直接返回
				return nil, err
			}
			if !ok {
				return nil, status.Errorf(codes.ResourceExhausted,
					"限流")
			}
		}

		return handler(ctx, req)
	}
}
