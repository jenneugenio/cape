package models

import (
	"fmt"
	"time"

	"github.com/manifoldco/go-base64"
)

type Config struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Setup bool `json:"setup"`

	// EncryptionKey is used to encrypt data in the system.
	// Specifically we're using envelope encryption which
	// can be read more about here
	// https://cloud.google.com/kms/docs/envelope-encryption.
	// Here it is encrypted and will be decrypted by the
	// root key.
	EncryptionKey *base64.Value `json:"encryption_key"`

	// AuthKeypair is encrypted using the root key, similar, to how the
	// EncryptionKey is encrypted.
	AuthKeypair *base64.Value `json:"auth_keypair"`
}

func (c *Config) Validate() error {
	if !c.Setup {
		return fmt.Errorf("config setup must be true")
	}

	if c.EncryptionKey == nil {
		return fmt.Errorf("an encryption key must be set")
	}

	if c.AuthKeypair == nil {
		return fmt.Errorf("auth keypair must be set")
	}

	return nil
}

// NewConfig returns a new Config
func NewConfig(encryptionKey *base64.Value, authKeypair *base64.Value) (*Config, error) {
	cfg := &Config{
		ID:            NewID(),
		CreatedAt:     now(),
		Setup:         true,
		EncryptionKey: encryptionKey,
		AuthKeypair:   authKeypair,
	}

	return cfg, cfg.Validate()
}
