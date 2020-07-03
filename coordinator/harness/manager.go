package harness

import (
	"context"
	"io/ioutil"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

const AdminEmail = "admin@cape.com"
const AdminName = "admin"
const AdminPassword = "iamtheadmin"

// User represents a user in the cape coordinator
type User struct {
	Client   *coordinator.Client
	User     *primitives.User
	Password primitives.Password
	Token    *base64.Value
}

// Service represents a service in the cape coordinator
type Service struct {
	ID    database.ID
	Token *auth.APIToken
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

	email, err := primitives.NewEmail(AdminEmail)
	if err != nil {
		return nil, err
	}

	name, err := primitives.NewName(AdminName)
	if err != nil {
		return nil, err
	}

	u, err := client.Setup(ctx, name, email, password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Client:   client,
		User:     u,
		Password: password,
	}

	session, err := client.EmailLogin(ctx, email, password)
	if err != nil {
		return nil, err
	}

	user.Token = session.Token
	m.Admin = user

	return client, nil
}

// CreateDeprecatedPolicy creates a policy on the coordinator!
func (m *Manager) CreatePolicy(ctx context.Context, policyPath string) error {
	data, err := ioutil.ReadFile(policyPath)
	if err != nil {
		return err
	}

	policy, err := primitives.ParsePolicy(data)
	if err != nil {
		return err
	}

	policy, err = m.Admin.Client.CreateDeprecatedPolicy(ctx, policy)
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
