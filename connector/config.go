package connector

import (
	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// Config is a configuration object for the Connector
type Config struct {
	InstanceID     primitives.Label
	Port           int
	Token          *auth.APIToken
	CoordinatorURL *primitives.URL
}

// Validate returns an error if the provided config is invalid
func (c *Config) Validate() error {
	if c.Port > 65535 || c.Port < 1 {
		return errors.New(InvalidConfigCause, "Port must be between 1-65535")
	}

	if c.Token == nil {
		return errors.New(InvalidConfigCause, "Missing token")
	}

	if err := c.Token.Validate(); err != nil {
		return err
	}

	if err := c.InstanceID.Validate(); err != nil {
		return err
	}

	if c.CoordinatorURL == nil {
		return errors.New(InvalidConfigCause, "Missing coordinator url")
	}

	return c.CoordinatorURL.Validate()
}

// GetPort satisfies the framework.Config interface
func (c *Config) GetPort() int {
	return c.Port
}

// GetInstanceID satisfies the framework.Config interface
func (c *Config) GetInstanceID() primitives.Label {
	return c.InstanceID
}
