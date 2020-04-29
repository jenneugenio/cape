package connector

import (
	"context"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zerolog "github.com/grpc-ecosystem/go-grpc-middleware/logging/zerolog"
	"github.com/justinas/alice"
	"github.com/manifoldco/healthz"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/capeprivacy/cape/auth"
	pb "github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/connector/sources"
	"github.com/capeprivacy/cape/framework"
	errors "github.com/capeprivacy/cape/partyerrors"
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
		logger:      logger,
	}

	grpcServer := grpc.NewServer(grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		errorStreamInterceptor,
		authStreamInterceptor,
		requestIDStreamInterceptor,
		grpc_zerolog.StreamServerInterceptor(logger, grpc_zerolog.WithCodes(handleCodes)),
	)))

	pb.RegisterDataConnectorServer(grpcServer, hndler)

	root := http.NewServeMux()
	root.Handle("/v1/version", framework.VersionHandler(cfg.InstanceID.String()))

	chain := alice.New(
		framework.RequestIDMiddleware,
		framework.LogMiddleware(logger),
		framework.RoundtripLoggerMiddleware,
		framework.RecoveryMiddleware,
		gziphandler.GzipHandler,
		cors.Default().Handler,
	).Then(healthz.NewHandler(root))

	return &Connector{
		cfg:         cfg,
		handler:     grpcHandlerFunc(grpcServer, chain),
		coordinator: coordinator,
		logger:      logger,
		cache:       cache,
	}, nil
}

func handleCodes(err error) codes.Code {
	if err == nil {
		return codes.OK
	} else if errors.CausedBy(err, auth.AuthorizationFailure) {
		return codes.PermissionDenied
	} else if errors.CausedBy(err, auth.InvalidAuthHeader) {
		return codes.Unauthenticated
	}

	return codes.Unknown
}
