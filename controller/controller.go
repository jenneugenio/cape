package controller

import (
	"context"
	"fmt"
	"net/url"

	"time"

	"github.com/dropoutlabs/privacyai/database"
)

// Controller is the central brain of PrivacyAI.  It keeps track of system
// users, policy, etc
type Controller struct {
	backend database.Backend
	name    string
}

// Start the controller
func (c *Controller) Start() {
	defer c.Stop()

	err := c.backend.Open(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(5 * time.Minute)
}

// Stop the controller
func (c *Controller) Stop() {
	c.backend.Close()
}

// New returns a pointer to a controller instance
func New(dbURL *url.URL, serviceID string) (*Controller, error) {
	name := fmt.Sprintf("cape-controller-%s", serviceID)
	backend, err := database.New(dbURL, name)

	if err != nil {
		return nil, err
	}

	return &Controller{
		backend: backend,
		name:    name,
	}, nil
}
