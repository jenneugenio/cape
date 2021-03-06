package harness

import (
	"context"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/models"
)

const AdminEmail = models.Email("admin@cape.com")
const AdminName = models.Name("admin")
const AdminPassword = "iamtheadmin"

// User represents a user in the cape coordinator
type User struct {
	Client   *coordinator.Client
	User     models.User
	Password models.Password
	Token    *base64.Value
}

// Manager represents an application state manager on-top of the Coordinator's
// harness. It's job is to provide convenience functions for setting up a
// coordinator's application state, managing users, and other utilities that
// make it write end-to-end integration tests.
type Manager struct {
	h     *Harness
	Admin *User
}

// Setup sets up the application state for the cluster (e.g. the `setup`
// mutation for the coordinator). This results in the creation of an Admin user.
//
// An authenticated client for the admin is returned.
func (m *Manager) Setup(ctx context.Context) (*coordinator.Client, error) {
	client, err := m.h.Client()
	if err != nil {
		return nil, err
	}

	password, err := models.NewPassword(AdminPassword)
	if err != nil {
		return nil, err
	}

	u := models.NewUser(AdminName, AdminEmail, models.Credentials{})

	user := &User{
		Client:   client,
		User:     u,
		Password: password,
	}

	session, err := client.EmailLogin(ctx, AdminEmail, password)
	if err != nil {
		return nil, err
	}

	user.User.ID = session.UserID
	user.Token = session.Token
	m.Admin = user

	return client, nil
}

// URL returns the url of the coordinator
func (m *Manager) URL() (*models.URL, error) {
	return m.h.URL()
}
