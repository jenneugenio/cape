package connector

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zerolog/ctxzr"
	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	pb "github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/connector/sources"
	fw "github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/policy"
	"github.com/capeprivacy/cape/primitives"
	"github.com/capeprivacy/cape/query"
	"github.com/capeprivacy/cape/version"
)

type grpcHandler struct {
	coordinator Coordinator
	cache       *sources.Cache
	logger      *zerolog.Logger
}

// Query implementation of the DataConnectorClient interface (see data_connector.pb.go)
func (g *grpcHandler) Query(req *pb.QueryRequest, server pb.DataConnector_QueryServer) error {
	err := g.handleQuery(req, server)
	if err != nil {
		return err
	}

	return nil
}

func (g *grpcHandler) handleQuery(req *pb.QueryRequest, server pb.DataConnector_QueryServer) error {
	session := fw.Session(server.Context())
	logger := ctxzr.Extract(server.Context()).Logger
	policies := session.Policies

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
	collection := primitives.Collection(query.Collection())
	schemas, err := source.Schema(context.Background(), collection)
	if err != nil {
		return err
	}

	if len(schemas) == 0 {
		logger.Info().Err(auth.ErrNoMatchingPolicies).Msgf("No schemas found for query %s", query.Collection())
		return auth.ErrNoMatchingPolicies
	}

	evaluator := policy.NewEvaluator(query, schemas[0], policies...)
	query, err = evaluator.Evaluate()
	if err != nil {
		return err
	}

	transforms := evaluator.Transforms()
	if len(transforms) > 0 {
		transformStream, err := NewTransformStream(server, schemas[0], transforms)
		if err != nil {
			return err
		}
		return source.Query(context.Background(), transformStream, query, req.GetLimit(), req.GetOffset())
	}

	return source.Query(context.Background(), server, query, req.GetLimit(), req.GetOffset())
}

// Schema implementation of the DataConnectorClient interface (see data_connector.pb.go)
func (g *grpcHandler) Schema(ctx context.Context, req *pb.SchemaRequest) (*pb.SchemaResponse, error) {
	dataSource, err := primitives.NewLabel(req.GetDataSource())
	if err != nil {
		return nil, err
	}

	source, err := g.cache.Get(ctx, dataSource)
	if err != nil {
		return nil, err
	}

	schemas, err := source.Schema(context.Background(), primitives.Collection(primitives.Star))
	if err != nil {
		return nil, err
	}

	return &pb.SchemaResponse{
		Schemas: schemas,
	}, nil
}

func (g *grpcHandler) Version(context.Context, *pb.VersionRequest) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{
		Version:   version.Version,
		BuildDate: version.BuildDate,
	}, nil
}

func (g *grpcHandler) GetCoordinator() Coordinator {
	return g.coordinator
}
