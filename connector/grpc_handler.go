package connector

import (
	"context"

	"github.com/dropoutlabs/cape/auth"
	pb "github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/primitives"
	"google.golang.org/grpc/metadata"
)

type grpcHandler struct {
	controller *Controller
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
	_, err := g.handleAuth(server.Context())
	if err != nil {
		return err
	}

	dataSource := req.GetDataSource()

	// TODO pull schema/data
	schema := &pb.Schema{DataSource: dataSource}
	r := &pb.Record{Schema: schema}

	err = server.Send(r)
	if err != nil {
		return err
	}

	// TODO loop over data and send more!!!

	return nil
}
