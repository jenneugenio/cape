package harness

import (
	"context"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/manifoldco/go-base64"
)

const AdminEmail = "admin@cape.com"
const AdminName = "admin"
const AdminPassword = "iamtheadmin"

// User represents a user in the cape controller
type User struct {
	Client   *controller.Client
	User     *primitives.User
	Password []byte
	Token    *base64.Value
}

// Service represents a service in the cape controller
type Service struct {
	Token *auth.APIToken
}

// Manager represents an application state manager on-top of the Controller's
// harness. It's job is to provide convenience functions for setting up a
// controller's application state, managing users, and other utilities that
// make it write end-to-end integration tests.
type Manager struct {
	h         *Harness
	Admin     *User
	Connector *Service
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

	email, err := primitives.NewEmail(AdminEmail)
	if err != nil {
		return nil, err
	}

	u, err := primitives.NewUser(AdminName, email, creds.Package())
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

	session, err := client.Login(ctx, email, pw)
	if err != nil {
		return nil, err
	}

	user.Token = session.Token
	m.Admin = user

	return client, nil
}

// CreateService creates a service on the controller with the given APIToken and URL
func (m *Manager) CreateService(ctx context.Context, apiToken *auth.APIToken, serviceURL *primitives.URL) error {
	creds, err := apiToken.Credentials()
	if err != nil {
		return err
	}

	service, err := primitives.NewService(apiToken.Email, primitives.DataConnectorServiceType, serviceURL, creds.Package())
	if err != nil {
		return err
	}

	_, err = m.Admin.Client.CreateService(ctx, service)
	if err != nil {
		return err
	}

	m.Connector = &Service{
		Token: apiToken,
	}

	return nil
}

func (m *Manager) URL() (*primitives.URL, error) {
	return m.h.URL()
}
