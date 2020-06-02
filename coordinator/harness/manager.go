package harness

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/client"
	"io/ioutil"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

const AdminEmail = "admin@cape.com"
const AdminName = "admin"
const AdminPassword = "iamtheadmin"

// User represents a user in the cape coordinator
type User struct {
	Client   *client.Client
	User     *primitives.User
	Password primitives.Password
	Token    *base64.Value
}

// Service represents a service in the cape coordinator
type Service struct {
	ID    database.ID
	Token *auth.APIToken
}

// Source represents a source on the cape coordinator
type Source struct {
	Label primitives.Label
}

// Manager represents an application state manager on-top of the Coordinator's
// harness. It's job is to provide convenience functions for setting up a
// coordinator's application state, managing users, and other utilities that
// make it write end-to-end integration tests.
type Manager struct {
	h          *Harness
	Admin      *User
	Connector  *Service
	TestSource *Source
}

// Setup sets up the application state for the cluster (e.g. the `setup`
// mutation for the coordinator). This results in the creation of an Admin user.
//
// An authenticated client for the admin is returned.
func (m *Manager) Setup(ctx context.Context) (*client.Client, error) {
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

// ReportSchema will log in as the provided worker (through the API token)
// and report the provided token as that worker
func (m *Manager) ReportSchema(ctx context.Context, token *auth.APIToken, sourceID database.ID, schema primitives.SchemaDefinition) error {
	c, err := m.h.Client()
	if err != nil {
		return err
	}

	_, err = c.TokenLogin(ctx, token)
	if err != nil {
		return err
	}

	return c.ReportSchema(ctx, sourceID, schema)
}

// Registers a worker and returns a token for that worker
func (m *Manager) CreateWorker(ctx context.Context) (*auth.APIToken, error) {
	workerEmail, err := primitives.NewEmail("worker@cape.com")
	if err != nil {
		return nil, err
	}

	service, err := primitives.NewService(workerEmail, primitives.WorkerServiceType, nil)
	if err != nil {
		return nil, err
	}

	worker, err := m.Admin.Client.CreateService(ctx, service)
	if err != nil {
		return nil, err
	}

	token, _, err := m.Admin.Client.CreateToken(ctx, worker)
	return token, err
}

// CreateSource creates a source on the coordinator
func (m *Manager) CreateSource(ctx context.Context, dbURL *primitives.DBURL, serviceID database.ID) error {
	sourceLabel := primitives.Label("test-source")
	_, err := m.Admin.Client.AddSource(ctx, sourceLabel, dbURL, &serviceID)
	if err != nil {
		return err
	}

	m.TestSource = &Source{
		Label: sourceLabel,
	}

	return nil
}

// CreateService creates a service on the coordinator with the given APIToken and URL
func (m *Manager) CreateService(ctx context.Context, email string, serviceURL *primitives.URL) error {
	e, err := primitives.NewEmail(email)
	if err != nil {
		return err
	}

	service, err := primitives.NewService(e, primitives.DataConnectorServiceType, serviceURL)
	if err != nil {
		return err
	}

	service, err = m.Admin.Client.CreateService(ctx, service)
	if err != nil {
		return err
	}

	apiToken, _, err := m.Admin.Client.CreateToken(ctx, service)
	if err != nil {
		return err
	}

	m.Connector = &Service{
		ID:    service.ID,
		Token: apiToken,
	}

	return nil
}

// CreatePolicy creates a policy on the coordinator!
func (m *Manager) CreatePolicy(ctx context.Context, policyPath string) error {
	data, err := ioutil.ReadFile(policyPath)
	if err != nil {
		return err
	}

	policy, err := primitives.ParsePolicy(data)
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
