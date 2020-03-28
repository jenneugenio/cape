package connector

import (
	"fmt"
	"time"

	"github.com/dropoutlabs/cape/auth"
)

// Connector is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Connector struct {
	InstanceID string
	Port       int
	Token      *auth.APIToken
}

// Start the connector
func (c *Connector) Start() error {
	time.Sleep(5 * time.Minute)

	return nil
}

// New returns a pointer to a controller instance
func New(cfg *Config) (*Connector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Connector{
		InstanceID: fmt.Sprintf("cape-connector-%s", cfg.InstanceID),
		Port:       cfg.Port,
		Token:      cfg.Token,
	}, nil
}
