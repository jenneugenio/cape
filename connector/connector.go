package connector

import "time"

// Connector is the central brain of PrivacyAI.  It keeps track of system
// users, policy, etc
type Connector struct {
}

// Start the connector
func (c *Connector) Start() {
	time.Sleep(5 * time.Minute)
}

// New returns a pointer to a controller instance
func New() *Connector {
	return &Connector{}
}
