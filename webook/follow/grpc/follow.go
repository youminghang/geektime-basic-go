package grpc

import (
	"context"
	followv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/follow/v1"
	"gitee.com/geekbang/basic-go/webook/follow/domain"
	"gitee.com/geekbang/basic-go/webook/follow/service"
	"google.golang.org/grpc"
)

type FollowServiceServer struct {
	followv1.UnimplementedFollowServiceServer
	svc service.FollowRelationService
}

func NewFollowRelationServiceServer(svc service.FollowRelationService) *FollowServiceServer {
	return &FollowServiceServer{
		svc: svc,
	}
}

func (f *FollowServiceServer) Register(server grpc.ServiceRegistrar) {
	followv1.RegisterFollowServiceServer(server, f)
}

func (f *FollowServiceServer) FollowRelationList(ctx context.Context, request *followv1.FollowRelationListRequest) (*followv1.FollowRelationListResponse, error) {
	relationList, err := f.svc.GetFollowee(ctx, request.Follower, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}
	res := make([]*followv1.FollowRelation, 0, len(relationList))
	for _, relation := range relationList {
		res = append(res, f.convertToView(relation))
	}
	return &followv1.FollowRelationListResponse{
		FollowRelations: res,
	}, nil
}

func (f *FollowServiceServer) FollowRelationInfo(ctx context.Context, request *followv1.FollowRelationInfoRequest) (*followv1.FollowRelationInfoResponse, error) {
	info, err := f.svc.FollowInfo(ctx, request.Follower, request.Followee)
	if err != nil {
		return nil, err
	}
	return &followv1.FollowRelationInfoResponse{
		FollowRelation: f.convertToView(info),
	}, nil
}

func (f *FollowServiceServer) AddFollowRelation(ctx context.Context, request *followv1.AddFollowRelationRequest) (*followv1.AddFollowRelationResponse, error) {
	err := f.svc.Follow(ctx, request.Follower, request.Followee)
	return &followv1.AddFollowRelationResponse{}, err
}

func (f *FollowServiceServer) convertToView(relation domain.FollowRelation) *followv1.FollowRelation {
	return &followv1.FollowRelation{
		Followee: relation.Followee,
		Follower: relation.Follower,
	}
}
