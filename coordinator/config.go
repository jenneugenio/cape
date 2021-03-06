package coordinator

import (
	"io/ioutil"

	"github.com/manifoldco/go-base64"
	"sigs.k8s.io/yaml"

	"github.com/capeprivacy/cape/coordinator/db/crypto"
	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Config represents the configuration that needs to be provided to the
// Coordinator.
type Config struct {
	Version    int          `json:"version"`
	DB         *DBConfig    `json:"db" envconfig:"DB_URL"`
	InstanceID models.Label `json:"instance_id,omitempty"`
	Port       int          `json:"port"`

	// RootKey is used to encrypt/decrypt EncryptionKey and should
	// be stored in a separate config file in a secret or
	// other secure location.
	RootKey string `json:"root_key"`

	// The kdf algorithm is not externally configurable (e.g. not available on
	// the configuration file) as it's only required to be configurable for
	// testing.
	//
	// In future when we support more than one production algorithm we can
	// expose this feature to customers.
	CredentialProducerAlg models.CredentialsAlgType `json:"-"`

	// CertFile contains a path to the coordinators Certificate file.
	CertFile string `json:"tls_cert,omitempty" envconfig:"TLS_CERT"`

	// KeyFile contains a path to the coordinators TLS private key.
	KeyFile string `json:"tls_key,omitempty" envconfig:"TLS_KEY"`

	// Cors specifies the configuration for serving (or disabling)
	// CORS headers
	Cors CorsConfig `json:"cors"`

	// Required if no admin user has been c
	User *UserConfig `json:"user,omitempty"`
}

type UserConfig struct {
	Name     models.Name  `json:"name"`
	Email    models.Email `json:"email"`
	Password string       `json:"password"`
}

type CorsConfig struct {
	Enable      bool     `json:"enable"`
	AllowOrigin []string `json:"allow_origin,omitempty"`
}

// DBConfig represent the database configuration
type DBConfig struct {
	Addr *models.DBURL `json:"addr"`
}

// Decode implements envconfig.Decoder for decoding
// environment variables
func (db *DBConfig) Decode(value string) error {
	addr, err := models.NewDBURL(value)
	if err != nil {
		return err
	}

	db.Addr = addr

	return nil
}

// Validate returns an error if the DBConfig is invalid
func (db *DBConfig) Validate() error {
	if db.Addr == nil {
		return errors.New(InvalidConfigCause, "A db address must be provided")
	}

	return db.Addr.Validate()
}

// GetPort returns the port and completes the framwork.Config interface
func (c *Config) GetPort() int {
	return c.Port
}

// GetInstanceID returns the instance id to satisfy the framework.Component
// interface
func (c *Config) GetInstanceID() models.Label {
	return c.InstanceID
}

// Validate returns an error if the config is invalid
func (c *Config) Validate() error {
	if c.Version != 1 {
		return errors.New(InvalidConfigCause, "Version must be 1")
	}

	if c.Port > 65535 || c.Port < 1 {
		return errors.New(InvalidConfigCause, "Port must be between 1-65335")
	}

	if c.DB == nil {
		return errors.New(InvalidConfigCause, "Missing db block")
	}

	if err := c.DB.Validate(); err != nil {
		return err
	}

	if c.RootKey == "" {
		return errors.New(InvalidConfigCause, "Missing root key")
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

func NewConfig(port int, dbURL *models.DBURL) (*Config, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Version: 1,
		Port:    port,
		DB: &DBConfig{
			Addr: dbURL,
		},
		RootKey: "base64key://" + base64.New(key[:]).String(),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
