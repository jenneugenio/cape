package connector

import (
	"context"

	"github.com/rs/zerolog"

	pb "github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/connector/sources"
	"github.com/capeprivacy/cape/framework"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/policy"
	"github.com/capeprivacy/cape/primitives"
	"github.com/capeprivacy/cape/query"
	"github.com/capeprivacy/cape/version"
)

type grpcHandler struct {
	coordinator *Coordinator
	cache       *sources.Cache
	logger      *zerolog.Logger
}

func (g *grpcHandler) Query(req *pb.QueryRequest, server pb.DataConnector_QueryServer) error {
	err := g.handleQuery(req, server)
	if err != nil {
		return err
	}

	return nil
}

func (g *grpcHandler) handleQuery(req *pb.QueryRequest, server pb.DataConnector_QueryServer) error {
	identity := fw.Identity(server.Context())

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

	return source.Query(context.Background(), server, query, req.GetLimit(), req.GetOffset())
}

func (g *grpcHandler) Version(context.Context, *pb.VersionRequest) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{
		Version:   version.Version,
		BuildDate: version.BuildDate,
	}, nil
}

func (g *grpcHandler) ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error) {
	return g.coordinator.ValidateToken(ctx, tokenStr)
}
