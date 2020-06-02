package worker

import (
	"github.com/rs/zerolog"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// Config is a configuration object for the Worker
type Config struct {
	DatabaseURL    *primitives.DBURL
	Token          *auth.APIToken
	CoordinatorURL *primitives.URL
	Logger         *zerolog.Logger
}

func (c *Config) Validate() error {
	if c.Token == nil {
		return errors.New(InvalidConfigCause, "A token must be supplied")
	}

	err := c.Token.Validate()
	if err != nil {
		return err
	}

	if c.DatabaseURL == nil {
		return errors.New(InvalidConfigCause, "A database url must be supplied")
	}

	err = c.DatabaseURL.Validate()
	if err != nil {
		return err
	}

	if c.CoordinatorURL == nil {
		return errors.New(InvalidConfigCause, "A coordinator url must be supplied")
	}

	if c.Logger == nil {
		return errors.New(InvalidConfigCause, "A logger must be provided")
	}

	return c.CoordinatorURL.Validate()
}

func NewConfig(token *auth.APIToken, dbURL *primitives.DBURL, coordinatorURL *primitives.URL, logger *zerolog.Logger) (*Config, error) {
	c := &Config{
		Token:          token,
		DatabaseURL:    dbURL,
		CoordinatorURL: coordinatorURL,
		Logger:         logger,
	}

	return c, c.Validate()
}
