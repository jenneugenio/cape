package harness

import (
	"os"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

var (
	// MissingConfig occurs when the config is missing a value
	MissingConfig = errors.NewCause(errors.BadRequestCategory, "missing_config")
)

// Config is the connector harness config
type Config struct {
	ControllerURL       *primitives.URL
	dbURL               *primitives.DBURL
	sourceMigrationsDir string
}

// Validate returns an error if the struct contains invalid configuration
func (c *Config) Validate() error {
	if c.sourceMigrationsDir == "" {
		return errors.New(MissingConfig, "Cannot run tests, missing 'CAPE_DB_SEED_MIGRATIONS' environment variable")
	}

	return nil
}

// NewConfig returns an instantiated version of the Controller Harness configuration
func NewConfig(controllerURL *primitives.URL) (*Config, error) {
	dbURL, err := primitives.NewDBURL(os.Getenv("CAPE_DB_URL"))
	if err != nil {
		return nil, err
	}

	c := &Config{
		ControllerURL:       controllerURL,
		sourceMigrationsDir: os.Getenv("CAPE_DB_SEED_MIGRATIONS"),
		dbURL:               dbURL,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}