package controller

import (
	"os"

	"github.com/dropoutlabs/privacyai/database"
)

// Controller is the central brain of PrivacyAI.  It keeps track of system
// users, policy, etc
type Controller struct {
	backend database.Backend
}

// Start the controller
func (c *Controller) Start() {
}

// New returns a pointer to a controller instance
func New() (*Controller, error) {
	addr := os.Getenv("DB_ADDR")
	backend, err := database.New(addr)

	if err != nil {
		return nil, err
	}

	return &Controller{
		backend: backend,
	}, nil
}
