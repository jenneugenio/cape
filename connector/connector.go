package connector

import (
	"context"
	"net/http"
	"strings"

	"github.com/justinas/alice"
	"github.com/manifoldco/healthz"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/connector/sources"
)

const connectorCertFile = "connector/certs/localhost.crt"
const connectorKeyFile = "connector/certs/localhost.key"

// Connector is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Connector struct {
	cfg         *Config
	handler     http.Handler
	coordinator *Coordinator
	cache       *sources.Cache
	logger      *zerolog.Logger
}

// Setup starts the connector
func (c *Connector) Setup(ctx context.Context) (http.Handler, error) {
	return c.handler, nil
}

// Teardown tears down the server
func (c *Connector) Teardown(ctx context.Context) error {
	return c.cache.Close(ctx)
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
func New(cfg *Config, logger *zerolog.Logger) (*Connector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	coordinator := NewCoordinator(cfg.Token, logger)

	sCfg := &sources.Config{
		InstanceID: cfg.InstanceID,
		Logger:     logger,
	}
	cache, err := sources.NewCache(sCfg, coordinator, nil)
	if err != nil {
		return nil, err
	}

	hndler := &grpcHandler{
		coordinator: coordinator,
		cache:       cache,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataConnectorServer(grpcServer, hndler)

	mux := http.NewServeMux()
	health := healthz.NewHandler(mux)

	chain := alice.New().Then(health)

	return &Connector{
		cfg:         cfg,
		handler:     grpcHandlerFunc(grpcServer, chain),
		coordinator: coordinator,
		logger:      logger,
		cache:       cache,
	}, nil
}
