package connector

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	gm "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/primitives"
)

type TestValidator struct{}

func (t *TestValidator) ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error) {
	return &primitives.User{}, nil
}

type TestStream struct {
	ctx context.Context
	md  metadata.MD
}

func (t *TestStream) SetHeader(md metadata.MD) error {
	return nil
}

func (t *TestStream) SendHeader(md metadata.MD) error {
	t.md = md
	return nil
}

func (t *TestStream) SetTrailer(md metadata.MD) {}

func (t *TestStream) Context() context.Context {
	return t.ctx
}

func (t *TestStream) SendMsg(m interface{}) error {
	return nil
}

func (t *TestStream) RecvMsg(m interface{}) error {
	return nil
}

func TestInterceptors(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("test auth", func(t *testing.T) {
		wasCalled := false

		ctx := context.Background()

		md := metadata.Pairs("authorization", "Bearer: acooltoken")
		ctx = metadata.NewIncomingContext(ctx, md)

		info := &grpc.StreamServerInfo{IsServerStream: true}
		err := authStreamInterceptor(&TestValidator{}, &TestStream{ctx: ctx}, info,
			func(srv interface{}, stream grpc.ServerStream) error {
				wasCalled = true
				return nil
			})

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(wasCalled).To(gm.BeTrue())
	})

	t.Run("test missing auth header", func(t *testing.T) {
		wasCalled := false
		ctx := context.Background()

		info := &grpc.StreamServerInfo{IsServerStream: true}
		err := authStreamInterceptor(&TestValidator{}, &TestStream{ctx: ctx}, info,
			func(srv interface{}, stream grpc.ServerStream) error {
				wasCalled = true
				return nil
			})

		gm.Expect(err).To(gm.Equal(auth.ErrorInvalidAuthHeader))
		gm.Expect(wasCalled).To(gm.BeFalse())
	})

	t.Run("test error handling", func(t *testing.T) {
		wasCalled := false
		ctx := context.Background()

		info := &grpc.StreamServerInfo{IsServerStream: true}
		err := errorStreamInterceptor(&TestValidator{}, &TestStream{ctx: ctx}, info,
			func(srv interface{}, stream grpc.ServerStream) error {
				wasCalled = true
				return auth.ErrorInvalidAuthHeader
			})
		gm.Expect(wasCalled).To(gm.BeTrue())
		gm.Expect(err).ToNot(gm.BeNil())
		stat := status.Convert(err)
		gm.Expect(stat.Message()).To(gm.Equal("Cape Error"))
	})

	t.Run("test request id", func(t *testing.T) {
		wasCalled := false
		ctx := context.Background()

		stream := &TestStream{ctx: ctx}
		info := &grpc.StreamServerInfo{IsServerStream: true}
		err := requestIDStreamInterceptor(&TestValidator{}, stream, info,
			func(srv interface{}, stream grpc.ServerStream) error {
				wasCalled = true
				return nil
			})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(wasCalled).To(gm.BeTrue())

		strs := stream.md.Get("x-request-id")
		gm.Expect(len(strs)).To(gm.Equal(1))

		_, err = uuid.Parse(strs[0])
		gm.Expect(err).To(gm.BeNil())
	})
}
