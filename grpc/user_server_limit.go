package grpc

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/pkg/ratelimit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LimiterUserServer struct {
	limiter ratelimit.Limiter
	UserServiceServer
}

func (l *LimiterUserServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	const keyPattern = "limiter:user:get_by_id:%d"
	key := fmt.Sprintf(keyPattern, req.GetId())
	ok, err := l.limiter.Limit(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(codes.ResourceExhausted,
			"限流")
	}
	return l.UserServiceServer.GetById(ctx, req)
}
