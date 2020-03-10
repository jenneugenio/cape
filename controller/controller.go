package controller

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/graph"
	"github.com/dropoutlabs/cape/graph/generated"
	"net/http"
	"net/url"
	"sync"
)

// Controller is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Controller struct {
	backend   database.Backend
	name      string
	gqlServer *http.Server
	wg        *sync.WaitGroup
}

func (c *Controller) startGQLServer() {
	go func() {
		defer c.wg.Done()

		if err := c.gqlServer.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

// Start the controller
func (c *Controller) Start(ctx context.Context) {
	err := c.backend.Open(ctx)
	if err != nil {
		fmt.Println(err)
	}

	c.wg.Add(1)
	c.startGQLServer()
}

// Stop the controller
func (c *Controller) Stop(ctx context.Context) error {
	defer c.backend.Close()
	err := c.gqlServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	c.wg.Wait()
	return nil
}

// New returns a pointer to a controller instance
func New(dbURL *url.URL, serviceID string) (*Controller, error) {
	name := fmt.Sprintf("cape-controller-%s", serviceID)
	backend, err := database.New(dbURL, name)

	if err != nil {
		return nil, err
	}

	gqlSrv := &http.Server{Addr: ":8081"}
	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Backend: backend,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", gqlHandler)

	httpServerExitDone := &sync.WaitGroup{}

	return &Controller{
		backend:   backend,
		name:      name,
		gqlServer: gqlSrv,
		wg:        httpServerExitDone,
	}, nil
}
