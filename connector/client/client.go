package client

import (
	"context"
	"crypto/x509"

	"github.com/manifoldco/go-base64"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/primitives"
)

// Client is a wrapper around the grpc client that
// connects to the connector and sends queries
type Client struct {
	client    pb.DataConnectorClient
	authToken *base64.Value
	conn      *grpc.ClientConn
}

// NewClient dials the connector and creates a client
func NewClient(connectorURL *primitives.URL, authToken *base64.Value, certPool *x509.CertPool) (*Client, error) {
	creds := credentials.NewClientTLSFromCert(certPool, "")

	// strip https from url, dial expects the protocol not be specified
	conn, err := grpc.Dial(connectorURL.String()[8:], grpc.WithBlock(),
		grpc.WithTransportCredentials(creds), grpc.WithStreamInterceptor(authClientInterceptor(authToken)))
	if err != nil {
		return nil, err
	}

	return &Client{
		client:    pb.NewDataConnectorClient(conn),
		conn:      conn,
		authToken: authToken,
	}, nil
}

// Query queries the data source with the specified query
func (c *Client) Query(ctx context.Context, dataSource primitives.Label, query string,
	limit int64, offset int64) (Stream, error) {
	req := &pb.QueryRequest{
		DataSource: dataSource.String(),
		Query:      query,
		Limit:      limit,
		Offset:     offset,
	}

	client, err := c.client.Query(ctx, req)
	if err != nil {
		return nil, err
	}

	stream := NewStream(ctx, client)

	return stream, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}
