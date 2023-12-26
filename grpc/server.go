package grpc

import (
	"context"
	"log"
)

type Server struct {
	UnimplementedUserServiceServer
	name string
}

func (s *Server) GetById(
	ctx context.Context,
	req *GetByIdReq) (*GetByIdResp, error) {
	log.Println("命中服务器", s.name)
	return &GetByIdResp{
		User: &User{
			Id:   req.Id,
			Name: s.name,
		},
	}, nil
}
