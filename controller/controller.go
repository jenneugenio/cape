package controller

// Controller is the central brain of PrivacyAI.  It keeps track of system
// users, policy, etc
type Controller struct {
}

// Start the controller
func (c *Controller) Start() {
}

// New returns a pointer to a controller instance
func New() *Controller {
	return &Controller{}
}
