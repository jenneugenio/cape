package connector

import (
	pb "github.com/dropoutlabs/cape/connector/proto"
)

type grpcHandler struct{}

func (g *grpcHandler) Query(req *pb.QueryRequest, server pb.DataConnector_QueryServer) error {
	dataSource := req.GetDataSource()

	// TODO pull schema/data
	schema := &pb.Schema{DataSource: dataSource}
	r := &pb.Record{Schema: schema}

	err := server.Send(r)
	if err != nil {
		return err
	}

	// TODO loop over data and send more!!!

	return nil
}
