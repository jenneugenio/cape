package connector

import (
	"context"
	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	coor "github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

type CoordinatorProvider interface {
	GetCoordinator() Coordinator
}

type Coordinator interface {
	ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error)
	GetIdentityPolicies(ctx context.Context, id database.ID) ([]*primitives.Policy, error)
	GetSourceByLabel(ctx context.Context, label primitives.Label, opts *coor.SourceOptions) (*coor.SourceResponse, error)
}

// Coordinator wraps a coordinator client giving some extra features
// such as lazily authenticated (i.e. only authenticating when necessary)
// and validating that a given token has a valid session
type coordinator struct {
	*coor.Client
	url    *primitives.URL
	token  *auth.APIToken
	logger *zerolog.Logger
}

// NewCoordinator returns a new Coordinator
func NewCoordinator(url *primitives.URL, token *auth.APIToken, logger *zerolog.Logger) (Coordinator, error) {
	transport, err := coor.NewReAuthTransport(url, token, logger)
	if err != nil {
		return nil, err
	}

	return &coordinator{
		Client: coor.NewClient(transport),
		token:  token,
		url:    url,
		logger: logger,
	}, nil
}

// ValidateToken validates that a given token has a valid session and returns the
// related identity. Used for validating that a user making a request actually
// has a valid session.
func (c *coordinator) ValidateToken(ctx context.Context, tokenStr string) (primitives.Identity, error) {
	token, err := auth.GetBearerToken(tokenStr)
	if err != nil {
		return nil, err
	}

	transport := coor.NewHTTPTransport(c.url, token)
	userClient := coor.NewClient(transport)

	identity, err := userClient.Me(ctx)
	if err != nil {
		c.logger.Info().Err(err).Msgf("Unable to validate token with coordinator at %s", c.url)
		return nil, err
	}

	return identity, nil
}
