package connector

import (
	"context"

	"github.com/manifoldco/go-base64"
	"google.golang.org/grpc"
)

func authClientInterceptor(authToken *base64.Value) grpc.StreamClientInterceptor { // nolint: deadcode,unused
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		creds := tokenAccess{token: authToken}
		opts = append(opts, grpc.PerRPCCredentials(creds))

		return streamer(ctx, desc, cc, method, opts...)
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
