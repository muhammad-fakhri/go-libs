package grpcmiddleware_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/muhammad-fakhri/go-libs/grpcmiddleware"
	"github.com/muhammad-fakhri/go-libs/grpcmiddleware/example/proto/example"
	pb "github.com/muhammad-fakhri/go-libs/grpcmiddleware/example/proto/example"
	"github.com/muhammad-fakhri/go-libs/grpcmiddleware/example/server"
	"github.com/muhammad-fakhri/go-libs/log"
	"google.golang.org/grpc"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMiddleware(t *testing.T) {
	go server.InitServer()
	//initialize client
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", server.ServerPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	userID := 1200704667
	c := pb.NewExampleServiceClient(conn)

	Convey("Success With Metadata", t, func() {
		reqctx := context.Background()
		d := log.CommonFields{
			ContextID: uuid.NewString(),
			Country:   "ID",
			UserID:    fmt.Sprintf("%d", userID),
			EventID:   "test event",
		}
		ctx := context.WithValue(reqctx, log.ContextDataMapKey, d.ToDataMap())
		ctx = grpcmiddleware.CreateRequestMetadata(ctx)

		_, err = c.Ping(ctx, &empty.Empty{})
		_, err = c.UserInfo(ctx, &example.UserInfoRequest{
			UserId: int64(userID),
			Base: &example.CommonRequest{
				TraceId: uuid.NewString(),
			},
		})
	})

	Convey("Success Without Metadata", t, func() {
		ctx := context.Background()
		_, err = c.Ping(ctx, &empty.Empty{})
		So(err, ShouldBeNil)

		resp, err := c.UserInfo(ctx, &example.UserInfoRequest{
			UserId: int64(userID),
			Base: &example.CommonRequest{
				TraceId: uuid.NewString(),
			},
		})
		So(err, ShouldBeNil)
		So(resp.Status, ShouldEqual, 1)
	})
}
