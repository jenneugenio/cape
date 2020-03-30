package connector

import (
	"github.com/dropoutlabs/cape/auth"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

// Config is a configuration object for the Connector
type Config struct {
	InstanceID primitives.Label
	Port       int
	Token      *auth.APIToken
}

// Validate returns an error if the provided config is invalid
func (c *Config) Validate() error {
	if c.Port > 65535 || c.Port < 1 {
		return errors.New(InvalidConfigCause, "Port must be between 1-65535")
	}

	if c.Token == nil {
		return errors.New(InvalidConfigCause, "Missing token")
	}

	return c.Token.Validate()
}

// GetPort satisfies the framework.Config interface
func (c *Config) GetPort() int {
	return c.Port
}
