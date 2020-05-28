package connector

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/client"
	fw "github.com/capeprivacy/cape/framework"
	"testing"

	"github.com/gofrs/uuid"
	gm "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

type testCoordinatorProvider struct {
	coordinator Coordinator
}

func (t *testCoordinatorProvider) GetCoordinator() Coordinator {
	return t.coordinator
}

type testCoordinator struct {
	cp auth.CredentialProducer
}

func newTestCoordinator() (*testCoordinator, error) {
	return &testCoordinator{
		cp: auth.DefaultSHA256Producer,
	}, nil
}

func (t *testCoordinator) ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error) {
	creds, err := t.cp.Generate("hellosignmeup")
	if err != nil {
		return nil, err
	}

	return primitives.NewUser("BOB GEORGE", primitives.Email{Email: "bob@george.com"}, creds)
}

func (t *testCoordinator) GetIdentityPolicies(ctx context.Context, id database.ID) ([]*primitives.Policy, error) {
	return []*primitives.Policy{}, nil
}

func (t *testCoordinator) GetSourceByLabel(ctx context.Context, label primitives.Label, opts *client.SourceOptions) (*client.SourceResponse, error) {
	return &client.SourceResponse{}, nil
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

	tc, err := newTestCoordinator()
	gm.Expect(err).To(gm.BeNil())

	cp := &testCoordinatorProvider{coordinator: tc}

	t.Run("test auth", func(t *testing.T) {
		wasCalled := false

		ctx := context.Background()
		md := metadata.Pairs("authorization", "Bearer: acooltoken")
		ctx = metadata.NewIncomingContext(ctx, md)

		info := &grpc.StreamServerInfo{IsServerStream: true}
		interceptor := Interceptor{cp}
		err := interceptor.AuthStreamInterceptor(interceptor.provider, &TestStream{ctx: ctx}, info,
			func(srv interface{}, stream grpc.ServerStream) error {
				wasCalled = true
				return nil
			})

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(wasCalled).To(gm.BeTrue())
	})

	t.Run("test auth unary", func(t *testing.T) {
		wasCalled := false

		ctx := context.Background()
		md := metadata.Pairs("authorization", "Bearer: acooltoken")
		ctx = metadata.NewIncomingContext(ctx, md)

		interceptor := Interceptor{cp}
		f := func(ctx context.Context, req interface{}) (interface{}, error) {
			wasCalled = true
			return nil, nil
		}

		var req interface{}
		_, err := interceptor.AuthUnaryInterceptor(ctx, req, nil, f)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(wasCalled).To(gm.BeTrue())
	})

	t.Run("test missing auth header stream", func(t *testing.T) {
		wasCalled := false
		ctx := context.Background()

		info := &grpc.StreamServerInfo{IsServerStream: true}
		interceptor := Interceptor{cp}
		err := interceptor.AuthStreamInterceptor(interceptor.provider, &TestStream{ctx: ctx}, info,
			func(srv interface{}, stream grpc.ServerStream) error {
				wasCalled = true
				return nil
			})

		gm.Expect(err).To(gm.Equal(auth.ErrorInvalidAuthHeader))
		gm.Expect(wasCalled).To(gm.BeFalse())
	})

	t.Run("test missing auth header unary", func(t *testing.T) {
		wasCalled := false

		ctx := context.Background()

		interceptor := Interceptor{cp}
		f := func(ctx context.Context, req interface{}) (interface{}, error) {
			wasCalled = true
			return nil, nil
		}

		var req interface{}
		_, err := interceptor.AuthUnaryInterceptor(ctx, req, nil, f)
		gm.Expect(err).To(gm.Equal(auth.ErrorInvalidAuthHeader))
		gm.Expect(wasCalled).To(gm.BeFalse())
	})

	t.Run("test error handling stream", func(t *testing.T) {
		wasCalled := false
		ctx := context.Background()

		info := &grpc.StreamServerInfo{IsServerStream: true}
		interceptor := Interceptor{cp}
		err := interceptor.ErrorStreamInterceptor(interceptor.provider, &TestStream{ctx: ctx}, info,
			func(srv interface{}, stream grpc.ServerStream) error {
				wasCalled = true
				return auth.ErrorInvalidAuthHeader
			})
		gm.Expect(wasCalled).To(gm.BeTrue())
		gm.Expect(err).ToNot(gm.BeNil())
		stat := status.Convert(err)
		gm.Expect(stat.Message()).To(gm.Equal("Cape Error"))
	})

	t.Run("test error handling unary", func(t *testing.T) {
		wasCalled := false

		ctx := context.Background()
		interceptor := Interceptor{cp}
		f := func(ctx context.Context, req interface{}) (interface{}, error) {
			wasCalled = true
			return nil, auth.ErrorInvalidAuthHeader
		}

		var req interface{}
		resp, err := interceptor.ErrorUnaryInterceptor(ctx, req, nil, f)
		gm.Expect(resp).To(gm.BeNil())
		gm.Expect(wasCalled).To(gm.BeTrue())
		gm.Expect(err).ToNot(gm.BeNil())
		stat := status.Convert(err)
		gm.Expect(stat.Message()).To(gm.Equal("Cape Error"))
	})

	t.Run("test request id stream", func(t *testing.T) {
		wasCalled := false
		ctx := context.Background()

		stream := &TestStream{ctx: ctx}
		info := &grpc.StreamServerInfo{IsServerStream: true}
		interceptor := Interceptor{cp}
		err := interceptor.RequestIDStreamInterceptor(interceptor.provider, stream, info,
			func(srv interface{}, stream grpc.ServerStream) error {
				wasCalled = true
				return nil
			})
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(wasCalled).To(gm.BeTrue())

		strs := stream.md.Get("x-request-id")
		gm.Expect(len(strs)).To(gm.Equal(1))

		_, err = uuid.FromString(strs[0])
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("test request id unary", func(t *testing.T) {
		wasCalled := false

		var tag interface{}

		ctx := context.Background()
		interceptor := Interceptor{cp}
		f := func(ctx context.Context, req interface{}) (interface{}, error) {
			wasCalled = true
			tag = ctx.Value(fw.RequestIDContextKey)
			return nil, nil
		}

		var req interface{}
		_, err := interceptor.RequestIDUnaryInterceptor(ctx, req, nil, f)

		gm.Expect(err).To(gm.BeNil())
		gm.Expect(wasCalled).To(gm.BeTrue())
		gm.Expect(tag).ToNot(gm.BeNil())
	})
}
