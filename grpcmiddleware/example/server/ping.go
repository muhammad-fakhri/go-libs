package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *GameServer) Ping(ctx context.Context, emp *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
