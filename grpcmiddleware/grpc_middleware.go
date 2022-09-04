package grpcmiddleware

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/muhammad-fakhri/go-libs/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	statusOK  = 0
	statusNOK = 1

	typeIngressGRPC = "ingress_grpc"
)

func PayloadUnaryServerLogInterceptor(logger log.SLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		rr := &log.RequestResponse{}

		rctx := getRequestContext(ctx)
		resp, err := handler(rctx, req)
		if err == nil {
			rr.Status = statusOK
		} else {
			rr.Status = statusNOK
		}

		if resp != nil {
			rr.ResponseBody = resp
		}

		rr.Type = typeIngressGRPC
		rr.RequestBody = req
		rr.URLPath = info.FullMethod
		rr.DurationMs = time.Since(startTime).Milliseconds()
		rr.RequestTimestamp = startTime

		if err != nil {
			logger.LogRequestResponse(rctx, rr, err)
		} else {
			logger.LogRequestResponse(rctx, rr)
		}

		return resp, err
	}
}

func getRequestContext(ctx context.Context) context.Context {
	data := make(map[string]string, 0)
	v, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	if u := v[log.ContextUserIdKey]; u != nil && len(u) > 0 {
		data[log.ContextUserIdKey] = u[0]
	}
	if c := v[log.ContextCountryKey]; c != nil && len(c) > 0 {
		data[log.ContextCountryKey] = c[0]
	}
	if e := v[log.ContextEventIdKey]; e != nil && len(e) > 0 {
		data[log.ContextEventIdKey] = e[0]
	}
	if i := v[log.ContextIdKey]; i != nil && len(i) > 0 {
		data[log.ContextIdKey] = i[0]
	} else {
		data[log.ContextIdKey] = uuid.New().String()
	}

	return context.WithValue(ctx, log.ContextDataMapKey, data)
}
