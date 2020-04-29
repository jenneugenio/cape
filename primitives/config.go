package primitives

import (
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type Config struct {
	*database.Primitive
	Setup bool `json:"setup"`
}

func (c *Config) Validate() error {
	if err := c.Primitive.Validate(); err != nil {
		return err
	}

	if !c.Setup {
		return errors.New(InvalidConfigCause, "Config setup must be true")
	}

	return nil
}

// GetType returns the type for this entity
func (c *Config) GetType() types.Type {
	return ConfigType
}

// NewConfig returns a new Config primitive
func NewConfig() (*Config, error) {
	p, err := database.NewPrimitive(ConfigType)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Primitive: p,
		Setup:     true,
	}

	return cfg, cfg.Validate()
}
