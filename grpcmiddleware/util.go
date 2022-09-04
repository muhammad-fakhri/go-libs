package grpcmiddleware

import (
	"context"

	"github.com/muhammad-fakhri/go-libs/log"
	"google.golang.org/grpc/metadata"
)

// client helper, to add common fields from context "value" (see go-libs/log) to metadata
func CreateRequestMetadata(ctx context.Context) context.Context {
	dataMap := ctx.Value(log.ContextDataMapKey)
	if dataMap == nil {
		return ctx
	}

	md := metadata.New(dataMap.(map[string]string))

	return metadata.NewOutgoingContext(ctx, md)
}
