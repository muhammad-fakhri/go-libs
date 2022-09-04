package main

import (
	"context"
	"fmt"

	"github.com/muhammad-fakhri/go-libs/grpcmiddleware"
	"github.com/muhammad-fakhri/go-libs/log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/muhammad-fakhri/go-libs/grpcmiddleware/example/proto/example"
	pb "github.com/muhammad-fakhri/go-libs/grpcmiddleware/example/proto/example"
	"github.com/muhammad-fakhri/go-libs/grpcmiddleware/example/server"
	"google.golang.org/grpc"
)

func main() {
	//initialize server
	go server.InitServer()

	//initialize client
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", server.ServerPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	userID := 1200704667
	c := pb.NewExampleServiceClient(conn)
	d := log.CommonFields{
		ContextID: uuid.NewString(),
		Country:   "ID",
		UserID:    fmt.Sprintf("%d", userID),
		EventID:   "test event",
	}

	reqctx := context.Background()
	ctx := context.WithValue(reqctx, log.ContextDataMapKey, d.ToDataMap())
	ctx = grpcmiddleware.CreateRequestMetadata(ctx)

	_, err = c.Ping(ctx, &empty.Empty{})
	_, err = c.UserInfo(ctx, &example.UserInfoRequest{
		UserId: int64(userID),
		Base: &example.CommonRequest{
			TraceId: uuid.NewString(),
		},
	})
}
