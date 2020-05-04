package connector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	spb "github.com/golang/protobuf/ptypes/struct"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
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

type CoordinatorProvider interface {
	GetCoordinator() Coordinator
}

func authStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	c := srv.(CoordinatorProvider)
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return auth.ErrorInvalidAuthHeader
	}

	authToken, ok := md["authorization"]
	if !ok {
		return auth.ErrorInvalidAuthHeader
	}

	coordinator := c.GetCoordinator()

	identity, err := coordinator.ValidateToken(ss.Context(), authToken[0])
	if err != nil {
		return err
	}

	policies, err := coordinator.GetIdentityPolicies(ss.Context(), identity.GetID())
	if err != nil {
		return err
	}

	session, err := auth.NewSession(identity, &primitives.Session{}, policies, []*primitives.Role{})
	if err != nil {
		return err
	}

	wStream := grpc_middleware.WrapServerStream(ss)
	wStream.WrappedContext = context.WithValue(wStream.WrappedContext, fw.SessionContextKey, session)

	return handler(srv, wStream)
}

func errorStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return returnGRPCError(handler(srv, ss))
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

func requestIDStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	id, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Sprintf("Could not generate a v4 uuid: %s", err))
	}

	wStream := grpc_middleware.WrapServerStream(ss)

	wStream.WrappedContext = context.WithValue(wStream.Context(), fw.RequestIDContextKey, id)

	md := metadata.Pairs("X-Request-ID", id.String())

	// this then gets scooped up by the logger interceptor
	tags := grpc_ctxtags.NewTags()
	tags.Set("request_id", id.String())
	wStream.WrappedContext = grpc_ctxtags.SetInContext(wStream.Context(), tags)

	err = wStream.SendHeader(md)
	if err != nil {
		return err
	}

	return handler(srv, wStream)
}
