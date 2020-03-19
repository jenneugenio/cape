package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/NYTimes/gziphandler"
	"github.com/justinas/alice"
	"github.com/rs/cors"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/manifoldco/healthz"
	"github.com/rs/zerolog"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/framework"
	"github.com/dropoutlabs/cape/graph"
	"github.com/dropoutlabs/cape/graph/generated"
)

// Controller is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Controller struct {
	backend    database.Backend
	instanceID string
	server     *http.Server
	logger     *zerolog.Logger
}

func (c *Controller) startGQLServer() error {
	err := c.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Start the controller
func (c *Controller) Start(ctx context.Context) error {
	err := c.backend.Open(ctx)
	if err != nil {
		fmt.Println(err)
	}

	c.logger.Info().Msg(fmt.Sprintf("Attempting to listen on: %s", c.server.Addr))
	return c.startGQLServer()
}

// Stop the controller
func (c *Controller) Stop(ctx context.Context) error {
	defer c.backend.Close()

	c.server.SetKeepAlivesEnabled(false)
	return c.server.Shutdown(ctx)
}

// New returns a pointer to a controller instance
func New(dbURL *url.URL, logger *zerolog.Logger, instanceID string) (*Controller, error) {
	backend, err := database.New(dbURL, instanceID)
	if err != nil {
		return nil, err
	}

	tokenAuth, err := auth.NewTokenAuthority(instanceID)
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
	root.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

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

	addr := ":8081"
	srv := &http.Server{
		Addr:    addr,
		Handler: chain,
	}

	return &Controller{
		backend:    backend,
		instanceID: instanceID,
		server:     srv,
		logger:     logger,
	}, nil
}
