package connector

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/primitives"
)

// Coordinator wraps a coordinator client giving some extra features
// such as lazily authenticated (i.e. only authenticating when necessary)
// and validating that a given token has a valid sessio
type Coordinator struct {
	*coordinator.Client
	token  *auth.APIToken
	logger *zerolog.Logger
}

// NewCoordinator returns a new Coordinator
func NewCoordinator(token *auth.APIToken, logger *zerolog.Logger) *Coordinator {
	return &Coordinator{
		token:  token,
		logger: logger,
	}
}

// ValidateToken validates that a given token has a valid session and returns the
// related identity. Used for validating that a user making a request actually
// has a valid session.
func (c *Coordinator) ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error) {
	// make sure the connector is actually authenticated before continuing
	err := c.authenticateClient(ctx)
	if err != nil {
		return nil, err
	}

	token, err := auth.GetBearerToken(tokenStr)
	if err != nil {
		return nil, err
	}

	transport := coordinator.NewTransport(c.token.URL, token)
	userClient := coordinator.NewClient(transport)

	return userClient.Me(ctx)
}

// authenticateClient lazily authenticates with a coordinator if required
func (c *Coordinator) authenticateClient(ctx context.Context) error {
	if c.Client == nil {
		transport := coordinator.NewTransport(c.token.URL, nil)
		c.Client = coordinator.NewClient(transport)
	}

	if c.Authenticated() {
		return nil
	}

	_, err := c.EmailLogin(ctx, c.token.Email, c.token.Secret)
	if err != nil {
		c.logger.Info().Msgf("Unable to log into the coordinator at %s. Err: %s", c.token.URL, err)
		return err
	}
	c.logger.Info().Msgf("Logged into the coordinator at %s", c.token.URL)

	return nil
}
