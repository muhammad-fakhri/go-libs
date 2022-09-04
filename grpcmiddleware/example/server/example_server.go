package server

import (
	"fmt"
	"net"

	"github.com/muhammad-fakhri/go-libs/grpcmiddleware"
	"github.com/muhammad-fakhri/go-libs/log"

	pb "github.com/muhammad-fakhri/go-libs/grpcmiddleware/example/proto/example"
	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	panichandler "github.com/kazegusuri/grpc-panic-handler"
)

const ServerPort = 7071

func InitServer() {
	addr := fmt.Sprintf(":%d", ServerPort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log := log.NewSLogger("test-middleware")
	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpcmiddleware.PayloadUnaryServerLogInterceptor(log),
			panichandler.UnaryPanicHandler),
	)

	pb.RegisterExampleServiceServer(s, &GameServer{Logger: log})

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
