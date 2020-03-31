package connector

import (
	"context"
	"crypto/x509"
	"fmt"
	"io"

	"github.com/manifoldco/go-base64"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/dropoutlabs/cape/connector/proto"
	"github.com/dropoutlabs/cape/primitives"
)

// Client is a wrapper around the grpc client that
// connects to the connector and sends queries
type Client struct {
	client    pb.DataConnectorClient
	authToken *base64.Value // nolint: unused, structcheck
}

// NewClient dials the connector and creates a client
func NewClient(connectorURL *primitives.URL, certPool *x509.CertPool) (*Client, error) {
	// TODO log in here
	creds := credentials.NewClientTLSFromCert(certPool, "")

	// strip https from url, dial expects the protocol not be specified
	conn, err := grpc.Dial(connectorURL.String()[8:], grpc.WithBlock(),
		grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: pb.NewDataConnectorClient(conn),
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
