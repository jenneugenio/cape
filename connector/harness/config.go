package harness

import (
	"os"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

var (
	// MissingConfig occurs when the config is missing a value
	MissingConfig = errors.NewCause(errors.BadRequestCategory, "missing_config")
)

// Config is the connector harness config
type Config struct {
	CoordinatorURL      *primitives.URL
	dbURL               *primitives.DBURL
	token               *auth.APIToken
	sourceMigrationsDir string
}

// Validate returns an error if the struct contains invalid configuration
func (c *Config) Validate() error {
	if c.sourceMigrationsDir == "" {
		return errors.New(MissingConfig, "Cannot run tests, missing 'CAPE_DB_SEED_MIGRATIONS' environment variable")
	}

	return nil
}

// NewConfig returns an instantiated version of the Coordinator Harness configuration
func NewConfig(coordinatorURL *primitives.URL, token *auth.APIToken) (*Config, error) {
	dbURL, err := primitives.NewDBURL(os.Getenv("CAPE_DB_URL"))
	if err != nil {
		return nil, err
	}

	c := &Config{
		CoordinatorURL:      coordinatorURL,
		sourceMigrationsDir: os.Getenv("CAPE_DB_SEED_MIGRATIONS"),
		token:               token,
		dbURL:               dbURL,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}
