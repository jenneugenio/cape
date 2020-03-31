package controller

import (
	"io/ioutil"

	"sigs.k8s.io/yaml"

	"github.com/dropoutlabs/cape/auth"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

// Config represents the configuration that needs to be provided to the
// Controller.
type Config struct {
	Version    int              `json:"version"`
	DB         *DBConfig        `json:"db"`
	InstanceID primitives.Label `json:"instance_id"`
	Port       int              `json:"port"`
	Auth       *AuthConfig      `json:"auth"`
}

// AuthConfig represents the authentication configuration
type AuthConfig struct {
	auth.KeypairPackage
}

// Validate returns an error if the auth config is invalid
func (a *AuthConfig) Validate() error {
	return a.KeypairPackage.Validate()
}

// Unpackage returns a Keypair for the given AuthConfig
func (a *AuthConfig) Unpackge() (*auth.Keypair, error) {
	return a.KeypairPackage.Unpackage()
}

//  DBConfig represent the database configuration
type DBConfig struct {
	Addr *primitives.DBURL `json:"addr"`
}

// Validate returns an error if the DBConfig is invalid
func (db *DBConfig) Validate() error {
	return db.Addr.Validate()
}

// GetPort returns the port and completes the framwork.Config interface
func (c *Config) GetPort() int {
	return c.Port
}

// GetInstanceID returns the instance id to satisfy the framework.Component
// interface
func (c *Config) GetInstanceID() primitives.Label {
	return c.InstanceID
}

// Validate returns an error if the config is invalid
func (c *Config) Validate() error {
	if c.Port > 65535 || c.Port < 1 {
		return errors.New(InvalidConfigCause, "Port must be between 1-65335")
	}

	if err := c.InstanceID.Validate(); err != nil {
		return err
	}

	if c.DB == nil {
		return errors.New(InvalidConfigCause, "Missing db block")
	}

	if err := c.DB.Validate(); err != nil {
		return err
	}

	return nil
}

// Write writes the configuration to the given filepath
func (c *Config) Write(filePath string) error {
	err := c.Validate()
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, b, 0700)
}

// LoadConfig parses a configuration file from given filepath and returns an
// initialized & validated config!
func LoadConfig(filePath string) (*Config, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
