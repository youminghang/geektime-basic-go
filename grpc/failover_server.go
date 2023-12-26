package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type FailoverServer struct {
	UnimplementedUserServiceServer
	code codes.Code
}

func (s *FailoverServer) GetById(
	ctx context.Context,
	req *GetByIdReq) (*GetByIdResp, error) {
	log.Println("命中了failover服务器")
	return &GetByIdResp{
		User: &User{
			Name: "failover 服务器",
		},
	}, status.Error(s.code, "模拟 failover")
}
