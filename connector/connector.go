package connector

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/justinas/alice"
	"github.com/manifoldco/healthz"
	"google.golang.org/grpc"

	"github.com/dropoutlabs/cape/auth"

	pb "github.com/dropoutlabs/cape/connector/proto"
)

const connectorCertFile = "connector/certs/localhost.crt"
const connectorKeyFile = "connector/certs/localhost.key"

// Connector is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Connector struct {
	InstanceID string
	Port       int
	Token      *auth.APIToken
	handler    http.Handler
}

// Setup starts the connector
func (c *Connector) Setup(ctx context.Context) (http.Handler, error) {
	return c.handler, nil
}

// Teardown tears down the server
func (c *Connector) Teardown(ctx context.Context) error {
	return nil
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(
			r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

// CertFiles returns the cert files used by ListenAndServeTLS
func (c *Connector) CertFiles() (certFile string, keyFile string) {
	certFile = connectorCertFile
	keyFile = connectorKeyFile
	return
}

// New returns a pointer to a Connector instance
func New(cfg *Config) (*Connector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataConnectorServer(grpcServer, &grpcHandler{})

	mux := http.NewServeMux()
	health := healthz.NewHandler(mux)

	chain := alice.New().Then(health)

	return &Connector{
		InstanceID: fmt.Sprintf("cape-connector-%s", cfg.InstanceID),
		Port:       cfg.Port,
		Token:      cfg.Token,
		handler:    grpcHandlerFunc(grpcServer, chain),
	}, nil
}
