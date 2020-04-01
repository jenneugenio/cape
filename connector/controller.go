package connector

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
)

// Controller wraps a controller client giving some extra features
// such as lazily authenticated (i.e. only authenticating when necessary)
// and validating that a given token has a valid sessio
type Controller struct {
	*controller.Client
	token  *auth.APIToken
	logger *zerolog.Logger
}

// NewController returns a new Controller
func NewController(token *auth.APIToken, logger *zerolog.Logger) *Controller {
	return &Controller{
		token:  token,
		logger: logger,
	}
}

// ValidateToken validates that a given token has a valid session and returns the
// related identity. Used for validating that a user making a request actually
// has a valid session.
func (c *Controller) ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error) {
	// make sure the connector is actually authenticated before continuing
	err := c.authenticateClient(ctx)
	if err != nil {
		return nil, err
	}

	token, err := auth.GetBearerToken(tokenStr)
	if err != nil {
		return nil, err
	}

	userClient := controller.NewClient(c.token.URL, token)

	return userClient.Me(ctx)
}

// authenticateClient lazily authenticates with a controller if required
func (c *Controller) authenticateClient(ctx context.Context) error {
	if c.Client == nil {
		c.Client = controller.NewClient(c.token.URL, nil)
	}

	if c.Authenticated() {
		return nil
	}

	_, err := c.Login(ctx, c.token.Email, c.token.Secret)
	if err != nil {
		c.logger.Info().Msgf("Unable to log into the controller at %s. Err: %s", c.token.URL, err)
		return err
	}
	c.logger.Info().Msgf("Logged into the controller at %s", c.token.URL)

	return nil
}
