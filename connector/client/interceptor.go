package client

import (
	"context"

	"github.com/manifoldco/go-base64"
	"google.golang.org/grpc"
)

func authClientStreamInterceptor(authToken *base64.Value) grpc.StreamClientInterceptor { // nolint: deadcode,unused
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if authToken != nil {
			creds := tokenAccess{token: authToken}
			opts = append(opts, grpc.PerRPCCredentials(creds))
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func authClientUnaryInterceptor(authToken *base64.Value) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if authToken != nil {
			creds := tokenAccess{token: authToken}
			opts = append(opts, grpc.PerRPCCredentials(creds))
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

type tokenAccess struct {
	token *base64.Value
}

func (t tokenAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token.String(),
	}, nil
}

func (t tokenAccess) RequireTransportSecurity() bool {
	return true
}
