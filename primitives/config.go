package primitives

import (
	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type Config struct {
	*database.Primitive
	Setup bool `json:"setup"`

	// EncryptionKey is used to encrypt data in the system.
	// Specifically we're using envelope encryption which
	// can be read more about here
	// https://cloud.google.com/kms/docs/envelope-encryption.
	// Here it is encrypted and will be decrypted by the
	// root key.
	EncryptionKey *base64.Value `json:"encryption_key"`

	Auth *base64.Value `json:"auth"`
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
func NewConfig(encryptionKey *base64.Value, auth *base64.Value) (*Config, error) {
	p, err := database.NewPrimitive(ConfigType)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Primitive:     p,
		Setup:         true,
		EncryptionKey: encryptionKey,
		Auth:          auth,
	}

	return cfg, cfg.Validate()
}

func (c *Config) GetEncryptable() bool {
	return false
}
