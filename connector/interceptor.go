package connector

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	spb "github.com/golang/protobuf/ptypes/struct"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/capeprivacy/cape/auth"
	fw "github.com/capeprivacy/cape/framework"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// Interceptor contains several interceptor methods that can be used in steam/unary chains
type Interceptor struct {
	provider CoordinatorProvider
}

func (i *Interceptor) AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	session, err := i.authIntercept(ctx)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, fw.SessionContextKey, session)
	return handler(ctx, req)
}

func (i *Interceptor) AuthStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	session, err := i.authIntercept(ss.Context())
	if err != nil {
		return err
	}

	wStream := grpc_middleware.WrapServerStream(ss)
	wStream.WrappedContext = context.WithValue(wStream.WrappedContext, fw.SessionContextKey, session)

	return handler(srv, wStream)
}

func (i *Interceptor) ErrorUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	err = returnGRPCError(err)

	return resp, err
}

func (i *Interceptor) ErrorStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return returnGRPCError(handler(srv, ss))
}

func (i *Interceptor) RequestIDUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	id, tags, err := i.requestIDIntercept()
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, fw.RequestIDContextKey, id)
	ctx = grpc_ctxtags.SetInContext(ctx, tags)

	return handler(ctx, req)
}

func (i *Interceptor) RequestIDStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	id, tags, err := i.requestIDIntercept()
	if err != nil {
		return err
	}

	wStream := grpc_middleware.WrapServerStream(ss)
	wStream.WrappedContext = context.WithValue(wStream.Context(), fw.RequestIDContextKey, id)
	md := metadata.Pairs("X-Request-ID", id.String())

	// tags gets scooped up by the logger interceptor
	wStream.WrappedContext = grpc_ctxtags.SetInContext(wStream.Context(), tags)
	err = wStream.SendHeader(md)
	if err != nil {
		return err
	}

	return handler(srv, wStream)
}


func (i *Interceptor) authIntercept(ctx context.Context) (*auth.Session, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, auth.ErrorInvalidAuthHeader
	}

	authToken, ok := md["authorization"]
	if !ok {
		return nil, auth.ErrorInvalidAuthHeader
	}

	coordinator := i.provider.GetCoordinator()
	identity, err := coordinator.ValidateToken(ctx, authToken[0])
	if err != nil {
		return nil, auth.ErrorInvalidAuthHeader
	}

	policies, err := coordinator.GetIdentityPolicies(ctx, identity.GetID())
	if err != nil {
		return nil, auth.ErrorInvalidAuthHeader
	}

	// TODO -- We aren't passing a credential provider here, but we don't actually need one
	// CredentialProvider is used by the coordinator to log a user in, we are using this object for policies
	return auth.NewSession(identity, &primitives.Session{}, policies, []*primitives.Role{}, nil)
}

func returnGRPCError(err error) error {
	if err == nil {
		return nil
	}

	pErr, ok := err.(*errors.Error)
	if !ok {
		pErr = errors.New(errors.UnknownCause, err.Error())
	}

	st := status.New(codes.Unknown, "Cape Error")

	details := &spb.Struct{}

	by, err := json.Marshal(pErr)
	if err != nil {
		return status.New(codes.Unknown, err.Error()).Err()
	}

	err = protojson.Unmarshal(by, details)
	if err != nil {
		return status.New(codes.Unknown, err.Error()).Err()
	}

	st, err = st.WithDetails(details)
	if err != nil {
		return status.New(codes.Unknown, err.Error()).Err()
	}

	return st.Err()
}

func (i *Interceptor) requestIDIntercept() (*uuid.UUID, grpc_ctxtags.Tags, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, nil, err
	}

	tags := grpc_ctxtags.NewTags()
	tags.Set("request_id", id.String())

	return &id, tags, nil
}