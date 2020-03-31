package primitives

import (
	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/database/types"
)

type Config struct {
	*database.Primitive
	Setup bool `json:"setup"`
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

	return &Config{
		Primitive: p,
		Setup:     true,
	}, nil
}
