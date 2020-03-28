package connector

import (
	"github.com/dropoutlabs/cape/auth"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

// Config is a configuration object for the Connector
type Config struct {
	InstanceID string
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
