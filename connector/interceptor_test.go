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
	coor "github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

type testCoordinatorProvider struct {
	coordinator Coordinator
}

func (t *testCoordinatorProvider) GetCoordinator() Coordinator {
	return t.coordinator
}

type testCoordinator struct{}

func (t *testCoordinator) ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error) {
	creds, _ := auth.NewCredentials([]byte("secret-hey"), nil)
	pCreds, _ := creds.Package()
	return primitives.NewUser("BOB GEORGE", primitives.Email{Email: "bob@george.com"}, pCreds)
}

func (t *testCoordinator) GetIdentityPolicies(ctx context.Context, id database.ID) ([]*primitives.Policy, error) {
	return []*primitives.Policy{}, nil
}

func (t *testCoordinator) GetSourceByLabel(ctx context.Context, label primitives.Label) (*coor.SourceResponse, error) {
	return &coor.SourceResponse{}, nil
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
		err := authStreamInterceptor(&testCoordinatorProvider{&testCoordinator{}}, &TestStream{ctx: ctx}, info,
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
		err := authStreamInterceptor(&testCoordinatorProvider{&testCoordinator{}}, &TestStream{ctx: ctx}, info,
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
		err := errorStreamInterceptor(&testCoordinatorProvider{&testCoordinator{}}, &TestStream{ctx: ctx}, info,
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
		err := requestIDStreamInterceptor(&testCoordinatorProvider{&testCoordinator{}}, stream, info,
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
}
