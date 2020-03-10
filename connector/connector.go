package connector

import (
	"fmt"
	"time"
)

// Connector is the central brain of Cape.  It keeps track of system
// users, policy, etc
type Connector struct {
	name string
}

// Start the connector
func (c *Connector) Start() {
	time.Sleep(5 * time.Minute)
}

// New returns a pointer to a controller instance
func New(serviceID string) *Connector {
	return &Connector{
		name: fmt.Sprintf("cape-connector-%s", serviceID),
	}
}
