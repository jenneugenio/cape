package connector

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/client"
	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

type CoordinatorProvider interface {
	GetCoordinator() Coordinator
}

type Coordinator interface {
	ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error)
	GetIdentityPolicies(ctx context.Context, id database.ID) ([]*primitives.Policy, error)
	GetSourceByLabel(ctx context.Context, label primitives.Label, opts *client.SourceOptions) (*client.SourceResponse, error)
}

// Coordinator wraps a coordinator client giving some extra features
// such as lazily authenticated (i.e. only authenticating when necessary)
// and validating that a given token has a valid session
type coordinator struct {
	*client.Client
	token  *auth.APIToken
	logger *zerolog.Logger
}

// NewCoordinator returns a new Coordinator
func NewCoordinator(token *auth.APIToken, logger *zerolog.Logger) Coordinator {
	return &coordinator{
		token:  token,
		logger: logger,
	}
}

// ValidateToken validates that a given token has a valid session and returns the
// related identity. Used for validating that a user making a request actually
// has a valid session.
func (c *coordinator) ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error) {
	// make sure the connector is actually authenticated before continuing
	//
	// XXX: This only ensures the connector is authenticated _IF_ you call
	// ValidateToken
	err := c.authenticateClient(ctx)
	if err != nil {
		return nil, err
	}

	token, err := auth.GetBearerToken(tokenStr)
	if err != nil {
		return nil, err
	}

	transport := client.NewTransport(c.token.URL, token)
	userClient := client.NewClient(transport)

	return userClient.Me(ctx)
}

// authenticateClient lazily authenticates with a coordinator if required
func (c *coordinator) authenticateClient(ctx context.Context) error {
	if c.Client == nil {
		transport := client.NewTransport(c.token.URL, nil)
		c.Client = client.NewClient(transport)
	}

	if c.Authenticated() {
		return nil
	}

	_, err := c.TokenLogin(ctx, c.token)
	if err != nil {
		c.logger.Info().Msgf("Unable to log into the coordinator at %s. Err: %s", c.token.URL, err)
		return err
	}
	c.logger.Info().Msgf("Logged into the coordinator at %s", c.token.URL)

	return nil
}
