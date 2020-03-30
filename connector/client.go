package connector

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/manifoldco/go-base64"
	"google.golang.org/grpc"

	pb "github.com/dropoutlabs/cape/connector/proto"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

// Client is a wrapper around the grpc client that
// connects to the connector and sends queries
type Client struct {
	client    pb.DataConnectorClient
	authToken *base64.Value
}

// NewClient dials the connector and creates a client
func NewClient(connectorURL *url.URL, authToken *base64.Value) (*Client, error) {
	if authToken == nil {
		return nil, errors.New(MissingAuthToken, "Must supply auth token when creating client")
	}

	conn, err := grpc.Dial(connectorURL.String(), grpc.WithInsecure(),
		grpc.WithBlock(), grpc.WithStreamInterceptor(authClientInterceptor(authToken)))
	if err != nil {
		return nil, err
	}

	return &Client{
		client:    pb.NewDataConnectorClient(conn),
		authToken: authToken,
	}, nil
}

// Query queries the datasource with the specified query
func (c *Client) Query(ctx context.Context, dataSource primitives.Label, query string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	req := &pb.QueryRequest{
		DataSource: dataSource.String(),
		Query:      query,
	}

	stream, err := c.client.Query(ctx, req)
	if err != nil {
		return err
	}

	i := 0
	var schema *pb.Schema
	for {
		record, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if i == 0 {
			schema = record.GetSchema()
		}
		fmt.Println(schema.GetFields())

		// TODO: handle record here
	}

	return nil
}
