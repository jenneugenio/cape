package connector

import (
	"bytes"
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/jsonpb"
	spb "github.com/golang/protobuf/ptypes/struct"

	"github.com/dropoutlabs/cape/auth"
	pb "github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/connector/sources"
	"github.com/dropoutlabs/cape/framework"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/policy"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/dropoutlabs/cape/query"
)

type grpcHandler struct {
	coordinator *Coordinator
	cache       *sources.Cache
}

func (g *grpcHandler) handleAuth(ctx context.Context) (primitives.Identity, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, auth.ErrorInvalidAuthHeader
	}

	authToken, ok := md["authorization"]
	if !ok {
		return nil, auth.ErrorInvalidAuthHeader
	}

	return g.coordinator.ValidateToken(ctx, authToken[0])
}

func (g *grpcHandler) Query(req *pb.QueryRequest, server pb.DataConnector_QueryServer) error {
	err := g.handleQuery(req, server)
	if err != nil {
		return returnGRPCError(err)
	}

	return nil
}

func (g *grpcHandler) handleQuery(req *pb.QueryRequest, server pb.DataConnector_QueryServer) error {
	identity, err := g.handleAuth(server.Context())
	if err != nil {
		return err
	}

	policies, err := g.coordinator.GetIdentityPolicies(server.Context(), identity.GetID())
	if err != nil {
		return err
	}

	if len(policies) == 0 {
		return framework.ErrAuthorization
	}

	dataSource, err := primitives.NewLabel(req.GetDataSource())
	if err != nil {
		return err
	}

	source, err := g.cache.Get(server.Context(), dataSource)
	if err != nil {
		return err
	}

	query, err := query.New(dataSource, req.GetQuery())
	if err != nil {
		return err
	}

	// This is using a new context because if its using the grpc context and
	// if the grpc connection is closed grpc cancels that context. This
	// causes pgx to ungracefully close the connection to postgres which
	// was causing a problems with k8s port forwarding causing our tests
	// to break.
	schema, err := source.Schema(context.Background(), query)
	if err != nil {
		return err
	}

	evaluator := policy.NewEvaluator(query, schema, policies...)
	query, err = evaluator.Evaluate()
	if err != nil {
		return err
	}

	return source.Query(context.Background(), query, server)
}

func returnGRPCError(err error) error {
	pErr, ok := err.(*errors.Error)
	if !ok {
		return status.New(codes.Unknown, err.Error()).Err()
	}

	st := status.New(codes.Unknown, "Cape Error")

	details := &spb.Struct{}

	by, err := json.Marshal(pErr)
	if err != nil {
		return status.New(codes.Unknown, err.Error()).Err()
	}

	err = jsonpb.Unmarshal(bytes.NewBuffer(by), details)
	if err != nil {
		return status.New(codes.Unknown, err.Error()).Err()
	}

	st, err = st.WithDetails(details)
	if err != nil {
		return status.New(codes.Unknown, err.Error()).Err()
	}

	return st.Err()
}
