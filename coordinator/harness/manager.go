package harness

import (
	"context"
	"io/ioutil"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

const AdminEmail = models.Email("admin@cape.com")
const AdminName = models.Name("admin")
const AdminPassword = "iamtheadmin"

// User represents a user in the cape coordinator
type User struct {
	Client   *coordinator.Client
	User     models.User
	Password primitives.Password
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

	password, err := primitives.NewPassword(AdminPassword)
	if err != nil {
		return nil, err
	}

	u := models.NewUser(AdminName, AdminEmail, nil)

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

// CreatePolicy creates a policy on the coordinator!
func (m *Manager) CreatePolicy(ctx context.Context, policyPath string) error {
	data, err := ioutil.ReadFile(policyPath)
	if err != nil {
		return err
	}

	policy, err := models.ParsePolicy(data)
	if err != nil {
		return err
	}

	policy, err = m.Admin.Client.CreatePolicy(ctx, policy)
	if err != nil {
		return err
	}

	role, err := m.Admin.Client.GetRoleByLabel(ctx, primitives.AdminRole)
	if err != nil {
		return err
	}

	_, err = m.Admin.Client.AttachPolicy(ctx, policy.ID, role.ID)
	if err != nil {
		return err
	}

	return nil
}

// URL returns the url of the coordinator
func (m *Manager) URL() (*primitives.URL, error) {
	return m.h.URL()
}
