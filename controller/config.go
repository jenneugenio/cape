package controller

import (
	"net/url"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

// Config represents the configuration that needs to be provided to the
// Controller.
type Config struct {
	DBURL      *url.URL
	InstanceID primitives.Label
	Port       int
}

// Validate returns an error if the config is invalid
func (c *Config) Validate() error {
	if c.Port > 65535 || c.Port < 1 {
		return errors.New(InvalidConfigCause, "Port must be between 1-65335")
	}

	if err := c.InstanceID.Validate(); err != nil {
		return err
	}

	return nil
}

// GetPort satisfies the framework.Config interface
func (c *Config) GetPort() int {
	return c.Port
}
