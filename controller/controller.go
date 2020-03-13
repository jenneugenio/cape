package controller

import (
	"github.com/dropoutlabs/cape/database"
)

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dropoutlabs/cape/graph"
	"github.com/dropoutlabs/cape/graph/generated"
	"net/http"
	"net/url"
)

// Controller is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Controller struct {
	backend   database.Backend
	name      string
	gqlServer *http.Server
}

func (c *Controller) startGQLServer() error {
	err := c.gqlServer.ListenAndServe()
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

	return c.startGQLServer()
}

// Stop the controller
func (c *Controller) Stop(ctx context.Context) error {
	defer c.backend.Close()

	c.gqlServer.SetKeepAlivesEnabled(false)
	return c.gqlServer.Shutdown(ctx)
}

// New returns a pointer to a controller instance
func New(dbURL *url.URL, serviceID string) (*Controller, error) {
	name := fmt.Sprintf("cape-controller-%s", serviceID)
	backend, err := database.New(dbURL, name)

	if err != nil {
		return nil, err
	}

	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Backend: backend,
	}}))

	mux := http.NewServeMux()

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", gqlHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	gqlSrv := &http.Server{Addr: ":8081", Handler: mux}

	return &Controller{
		backend:   backend,
		name:      name,
		gqlServer: gqlSrv,
	}, nil
}
