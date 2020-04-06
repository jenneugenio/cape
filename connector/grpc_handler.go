package connector

import (
	"context"

	"github.com/dropoutlabs/cape/auth"
	pb "github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/connector/sources"
	"github.com/dropoutlabs/cape/framework"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/dropoutlabs/cape/query"
	"google.golang.org/grpc/metadata"
)

type grpcHandler struct {
	controller *Controller
	cache      *sources.Cache
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

	return g.controller.ValidateToken(ctx, authToken[0])
}

func (g *grpcHandler) Query(req *pb.QueryRequest, server pb.DataConnector_QueryServer) error {
	identity, err := g.handleAuth(server.Context())
	if err != nil {
		return err
	}

	policies, err := g.controller.GetIdentityPolicies(server.Context(), identity.GetID())
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

	query, err = query.Rewrite(policies[0])
	if err != nil {
		return err
	}

	// This is using a new context because if its using the grpc context and
	// if the grpc connection is closed grpc cancels that context. This
	// causes pgx to ungracefully close the connection to postgres which
	// was causing a problems with k8s port forwarding causing our tests
	// to break.
	err = source.Query(context.Background(), query, server)
	if err != nil {
		return err
	}

	return nil
}
