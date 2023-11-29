package grpc

import (
	"context"
	"log"
)

type Server struct {
	UnimplementedUserServiceServer
}

func (s *Server) GetById(
	ctx context.Context,
	req *GetByIdReq) (*GetByIdResp, error) {
	log.Println(req)
	return &GetByIdResp{
		User: &User{
			Id:   123,
			Name: "测试用户",
		},
	}, nil
}
