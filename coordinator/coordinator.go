package coordinator

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NYTimes/gziphandler"
	"github.com/justinas/alice"
	"github.com/manifoldco/healthz"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/coordinator/graph"
	"github.com/dropoutlabs/cape/coordinator/graph/generated"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/framework"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

// Coordinator is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Coordinator struct {
	cfg     *Config
	backend database.Backend
	handler http.Handler
	logger  *zerolog.Logger
}

// Setup the coordinator so it's ready to be served!
func (c *Coordinator) Setup(ctx context.Context) (http.Handler, error) {
	err := c.backend.Open(ctx)
	if err != nil {
		return nil, err
	}

	return c.handler, nil
}

// Teardown the coordinator taking it back to it's start state!
func (c *Coordinator) Teardown(ctx context.Context) error {
	return c.backend.Close()
}

// CertFiles implements the Component interface. Coordinator doesn't support
// TLS right now so not needed!
func (c *Coordinator) CertFiles() (certFile string, keyFile string) {
	return
}

// New validates the input and returns a constructed Coordinator
func New(cfg *Config, logger *zerolog.Logger) (*Coordinator, error) {
	if cfg == nil {
		return nil, errors.New(InvalidConfigCause, "Config must be provided")
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	backend, err := database.New(cfg.DB.Addr.ToURL(), cfg.InstanceID.String())
	if err != nil {
		return nil, err
	}

	keypair, err := cfg.Auth.Unpackage()
	if err != nil {
		return nil, err
	}

	tokenAuth, err := auth.NewTokenAuthority(keypair, cfg.InstanceID.String())
	if err != nil {
		return nil, err
	}

	config := generated.Config{Resolvers: &graph.Resolver{
		Backend:        backend,
		TokenAuthority: tokenAuth,
	}}

	config.Directives.IsAuthenticated = framework.IsAuthenticatedDirective(backend, tokenAuth)
	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	root := http.NewServeMux()
	root.Handle("/v1", playground.Handler("GraphQL playground", "/query"))
	root.Handle("/v1/query", gqlHandler)
	root.Handle("/v1/version", framework.VersionHandler(cfg.InstanceID.String()))

	health := healthz.NewHandler(root)
	chain := alice.New(
		framework.RequestIDMiddleware,
		framework.LogMiddleware(logger),
		framework.AuthTokenMiddleware,
		framework.RoundtripLoggerMiddleware,
		framework.RecoveryMiddleware,
		gziphandler.GzipHandler,
		cors.Default().Handler,
	).Then(health)

	return &Coordinator{
		cfg:     cfg,
		handler: chain,
		backend: backend,
		logger:  logger,
	}, nil
}
