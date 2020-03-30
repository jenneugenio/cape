package connector

import (
	"context"

	"github.com/manifoldco/go-base64"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func authServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return nil
	}
}

func authClientInterceptor(authToken *base64.Value) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "token", authToken.String())
		return streamer(ctx, desc, cc, method, opts...)
	}
}
