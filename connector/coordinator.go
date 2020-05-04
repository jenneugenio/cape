package connector

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	coor "github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

type Coordinator interface {
	ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error)
	GetIdentityPolicies(ctx context.Context, id database.ID) ([]*primitives.Policy, error)
	GetSourceByLabel(ctx context.Context, label primitives.Label) (*primitives.Source, error)
}

// Coordinator wraps a coordinator client giving some extra features
// such as lazily authenticated (i.e. only authenticating when necessary)
// and validating that a given token has a valid sessio
type coordinator struct {
	*coor.Client
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
	err := c.authenticateClient(ctx)
	if err != nil {
		return nil, err
	}

	token, err := auth.GetBearerToken(tokenStr)
	if err != nil {
		return nil, err
	}

	transport := coor.NewTransport(c.token.URL, token)
	userClient := coor.NewClient(transport)

	return userClient.Me(ctx)
}

// authenticateClient lazily authenticates with a coordinator if required
func (c *coordinator) authenticateClient(ctx context.Context) error {
	if c.Client == nil {
		transport := coor.NewTransport(c.token.URL, nil)
		c.Client = coor.NewClient(transport)
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
