package worker

import (
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/primitives"
	"github.com/rs/zerolog"
)

// Config is a configuration object for the Worker
type Config struct {
	DatabaseURL *primitives.DBURL
	Token       *auth.APIToken
	Logger      *zerolog.Logger
}

func (c *Config) Validate() error {
	err := c.Token.Validate()
	if err != nil {
		return err
	}

	err = c.DatabaseURL.Validate()
	if err != nil {
		return err
	}

	return nil
}

func NewConfig(token *auth.APIToken, dbURL *primitives.DBURL, logger *zerolog.Logger) (*Config, error) {
	c := &Config{
		Token:       token,
		DatabaseURL: dbURL,
		Logger:      logger,
	}

	return c, c.Validate()
}
