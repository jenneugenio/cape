package worker

import (
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/primitives"
)

// Config is a configuration object for the Worker
type Config struct {
	DatabaseURL *primitives.DBURL
	Token       *auth.APIToken
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

func NewConfig(token *auth.APIToken, dbURL *primitives.DBURL) (*Config, error) {
	c := &Config{
		Token: token,
		DatabaseURL: dbURL,
	}

	return c, c.Validate()
}