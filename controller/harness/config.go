package harness

import (
	"os"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

var (
	MissingConfig = errors.NewCause(errors.BadRequestCategory, "missing_config")
)

// Config represents the configuration required to create a Controller Harness
type Config struct {
	dbURL             string
	migrationsDir     string
	migrationsTestDir string
}

// Validate returns an error if the struct contains invalid configuration
func (c *Config) Validate() error {
	if c.dbURL == "" {
		return errors.New(MissingConfig, "Cannot run tests, missing 'CAPE_DB_URL' environment variable")
	}

	if c.migrationsDir == "" {
		return errors.New(MissingConfig, "Cannot run tests, missing 'CAPE_DB_MIGRATIONS' environment variable")
	}

	if c.migrationsTestDir == "" {
		return errors.New(MissingConfig, "Cannot run tests, missing 'CAPE_DB_TEST_MIGRATIONS' environment variable")
	}

	return nil
}

// Migrations returns a list of directories containing migrations
func (c *Config) Migrations() []string {
	return []string{c.migrationsDir, c.migrationsTestDir}
}

// NewConfig returns an instantiated version of the Controller Harness configuration
func NewConfig() (*Config, error) {
	c := &Config{
		dbURL:             os.Getenv("CAPE_DB_URL"),
		migrationsDir:     os.Getenv("CAPE_DB_MIGRATIONS"),
		migrationsTestDir: os.Getenv("CAPE_DB_TEST_MIGRATIONS"),
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}
