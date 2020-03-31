package connector

import (
	"context"
	"fmt"
	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/justinas/alice"
	"github.com/manifoldco/healthz"
	"net/http"
)

// Connector is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Connector struct {
	InstanceID string
	Port       int
	Token      *auth.APIToken
	handler    http.Handler
}

func (c *Connector) login(ctx context.Context) (*controller.Client, error) {
	client := controller.NewClient(c.Token.URL, nil)
	_, err := client.Login(ctx, c.Token.Email, c.Token.Secret)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Start the connector
func (c *Connector) Setup(ctx context.Context) (http.Handler, error) {
	c.login(ctx)

	return c.handler, nil
}

func (c *Connector) Teardown(ctx context.Context) error {
	return nil
}

// New returns a pointer to a controller instance
func New(cfg *Config) (*Connector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	health := healthz.NewHandler(mux)
	chain := alice.New().Then(health)

	return &Connector{
		InstanceID: fmt.Sprintf("cape-connector-%s", cfg.InstanceID),
		Port:       cfg.Port,
		Token:      cfg.Token,
		handler:    chain,
	}, nil
}
