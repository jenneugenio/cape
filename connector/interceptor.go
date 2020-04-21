package connector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
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

type TokenValidator interface {
	ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error)
}

func authStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	t := srv.(TokenValidator)

	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		fmt.Println(":D", ok)
		return auth.ErrorInvalidAuthHeader
	}

	authToken, ok := md["authorization"]
	if !ok {
		return auth.ErrorInvalidAuthHeader
	}

	identity, err := t.ValidateToken(ss.Context(), authToken[0])
	if err != nil {
		return err
	}

	wStream := grpc_middleware.WrapServerStream(ss)
	wStream.WrappedContext = context.WithValue(wStream.WrappedContext, fw.IdentityContextKey, identity)

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

	// TODO make sure we use this ID on the logger as well
	wStream.WrappedContext = context.WithValue(wStream.Context(), fw.RequestIDContextKey, id)

	md := metadata.Pairs("X-Request-ID", id.String())

	err = wStream.SendHeader(md)
	if err != nil {
		return err
	}

	return handler(srv, wStream)
}
