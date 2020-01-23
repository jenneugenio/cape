package connector

// Connector is the central brain of PrivacyAI.  It keeps track of system
// users, policy, etc
type Connector struct {
}

// Start the connector
func (c *Connector) Start() {
}

// New returns a pointer to a controller instance
func New() *Connector {
	return &Connector{}
}
