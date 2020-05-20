package connector

import (
	"context"
	"github.com/rs/zerolog"

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
	schema, err := source.QuerySchema(context.Background(), query)
	if err != nil {
		return err
	}

	evaluator := policy.NewEvaluator(query, schema, policies...)
	query, err = evaluator.Evaluate()
	if err != nil {
		return err
	}

	transforms := evaluator.Transforms()
	if len(transforms) > 0 {
		transformStream, err := NewTransformStream(server, schema, transforms)
		if err != nil {
			return err
		}
		return source.Query(context.Background(), transformStream, query, req.GetLimit(), req.GetOffset())
	}

	return source.Query(context.Background(), server, query, req.GetLimit(), req.GetOffset())
}

// QuerySchema implementation of the DataConnectorClient interface (see data_connector.pb.go)
func (g *grpcHandler) Schema(ctx context.Context, req *pb.SchemaRequest) (*pb.SchemaResponse, error) {
	dataSource, err := primitives.NewLabel(req.GetDataSource())
	if err != nil {
		return nil, err
	}

	source, err := g.cache.Get(ctx, dataSource)
	if err != nil {
		return nil, err
	}

	// This is using a new context because if its using the grpc context and
	// if the grpc connection is closed grpc cancels that context. This
	// causes pgx to ungracefully close the connection to postgres which
	// was causing a problems with k8s port forwarding causing our tests
	// to break.
	schemas, err := source.SourceSchema(context.Background())
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
