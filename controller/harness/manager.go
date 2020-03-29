package harness

import (
	"context"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
)

const AdminEmail = "admin@cape.com"
const AdminName = "admin"
const AdminPassword = "iamtheadmin"

// User represents a user in the cape controller
type User struct {
	Client   *controller.Client
	User     *primitives.User
	Password []byte
}

// Manager represents an application state manager on-top of the Controller's
// harness. It's job is to provide convenience functions for setting up a
// controller's application state, managing users, and other utilities that
// make it write end-to-end integration tests.
type Manager struct {
	h     *Harness
	Admin *User
}

// Setup sets up the application state for the cluster (e.g. the `setup`
// mutation for the controller). This results in the creation of an Admin user.
//
// An authenticated client for the admin is returned.
func (m *Manager) Setup(ctx context.Context) (*controller.Client, error) {
	client, err := m.h.Client()
	if err != nil {
		return nil, err
	}

	pw := []byte(AdminPassword)
	creds, err := auth.NewCredentials(pw, nil)
	if err != nil {
		return nil, err
	}

	u, err := primitives.NewUser(AdminName, AdminEmail, creds.Package())
	if err != nil {
		return nil, err
	}

	u, err = client.Setup(ctx, u)
	if err != nil {
		return nil, err
	}

	user := &User{
		Client:   client,
		User:     u,
		Password: pw,
	}

	_, err = client.Login(ctx, "admin@cape.com", pw)
	if err != nil {
		return nil, err
	}

	m.Admin = user
	return client, nil
}
