package server

import (
	"context"

	pb "github.com/muhammad-fakhri/go-libs/grpcmiddleware/example/proto/example"
)

func (s *GameServer) UserInfo(ctx context.Context, req *pb.UserInfoRequest) (resp *pb.UserInfoResponse, err error) {
	status := int32(1)
	s.Logger.Infof(ctx, "get user status from data source, result: %d", status)
	return &pb.UserInfoResponse{
		Status: status,
	}, err
}
