package controller

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/dropoutlabs/cape/database"
	"github.com/dropoutlabs/cape/graph/model"

	"github.com/machinebox/graphql"
	"github.com/manifoldco/go-base64"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/primitives"
)

// Client is a wrapper around the graphql client that
// connects to the controller and sends queries
type Client struct {
	client    *graphql.Client
	authToken *base64.Value
}

// NewClient returns a new client that connects to the given
// url and set the struct required struct members
func NewClient(controllerURL *url.URL, authToken *base64.Value) *Client {
	return &Client{
		client:    graphql.NewClient(controllerURL.String() + "/v1/query"),
		authToken: authToken,
	}
}

// Raw wraps the NewRequest and does common req changes like adding authorization
// headers. It calls Run passing the object to be filled with the request data.
func (c *Client) Raw(ctx context.Context, query string,
	variables map[string]interface{}, resp interface{}) error {
	req := graphql.NewRequest(query)

	for key, val := range variables {
		req.Var(key, val)
	}

	if c.authToken != nil {
		req.Header.Add("Authorization", "Bearer "+c.authToken.String())
	}

	err := c.client.Run(ctx, req, resp)
	if err != nil {
		return err
	}

	return nil
}

// createLoginSession runs a CreateLoginSession mutation that creates a
// login session and returns it and also sets it on the
func (c *Client) createLoginSession(ctx context.Context, email primitives.Email) (*primitives.Session, error) {
	var resp struct {
		Session primitives.Session `json:"createLoginSession"`
	}

	variables := make(map[string]interface{})
	variables["email"] = email

	err := c.Raw(ctx, `
		mutation CreateLoginSession($email: Email!) {
			createLoginSession(input: { email: $email }) {
				id
				identity_id
				expires_at
				type
				token
				credentials {
					salt
					alg
				}
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Session, nil
}

// createAuthSession creates a authenticated session which can then be used
// for all other graphql queries. Replaces the old session set on the client
// and returns it
func (c *Client) createAuthSession(ctx context.Context, sig *base64.Value) (*primitives.Session, error) {
	var resp struct {
		Session primitives.Session `json:"createAuthSession"`
	}

	variables := make(map[string]interface{})
	variables["signature"] = sig

	err := c.Raw(ctx, `
		mutation CreateAuthSession($signature: Base64!) {
			createAuthSession(input: { signature: $signature }) {
				id
				identity_id
				expires_at
				type
				token
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Session, nil
}

// CreateUser creates a user and returns it
func (c *Client) CreateUser(ctx context.Context, user *primitives.User) (*primitives.User, error) {
	var resp struct {
		User primitives.User `json:"createUser"`
	}

	variables := make(map[string]interface{})
	variables["name"] = user.Name
	variables["email"] = user.Email
	variables["public_key"] = user.Credentials.PublicKey
	variables["salt"] = user.Credentials.Salt
	variables["alg"] = "EDDSA"

	err := c.Raw(ctx, `
		mutation CreateUser ($name: Name!, $email: Email!, $public_key: Base64!, $salt: Base64!, $alg: CredentialsAlgType!) {
			createUser(input: { name: $name, email: $email, public_key: $public_key, salt: $salt, alg: $alg}) {
				id
				name
				email
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.User, nil
}

// Login calls the CreateLoginSession and CreateAuthSession
// mutations
func (c *Client) Login(ctx context.Context, email primitives.Email, password []byte) (*primitives.Session, error) {
	session, err := c.createLoginSession(ctx, email)
	if err != nil {
		return nil, err
	}

	c.authToken = session.Token

	creds, err := auth.NewCredentials(password, session.Credentials.Salt)
	if err != nil {
		return nil, err
	}

	sig, err := creds.Sign(c.authToken)
	if err != nil {
		return nil, err
	}

	session, err = c.createAuthSession(ctx, sig)
	if err != nil {
		return nil, err
	}

	c.authToken = session.Token

	return session, nil
}

// Logout calls the deleteSession mutation
func (c *Client) Logout(ctx context.Context, authToken *base64.Value) error {
	token := authToken
	if authToken == nil {
		token = c.authToken
	}

	variables := make(map[string]interface{})
	variables["token"] = token

	return c.Raw(ctx, `
		mutation DeleteSession($token: Base64!) {
			deleteSession(input: { token: $token })
		}
	`, variables, nil)
}

// Role Routes

// CreateRole creates a new role with a label
func (c *Client) CreateRole(ctx context.Context, label primitives.Label, identityIDs []database.ID) (*primitives.Role, error) {
	var resp struct {
		Role primitives.Role `json:"createRole"`
	}

	variables := make(map[string]interface{})
	variables["ids"] = identityIDs
	variables["label"] = label

	err := c.Raw(ctx, `
		mutation CreateRole($label: Label!, $ids: [ID!]) {
			createRole(input: { label: $label, identity_ids: $ids }) {
				id
				label
				system
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Role, nil
}

// DeleteRole deletes a role with the given id
func (c *Client) DeleteRole(ctx context.Context, id database.ID) error {
	variables := make(map[string]interface{})
	variables["id"] = id

	return c.Raw(ctx, `
		mutation DeleteRole($id: ID!) {
			deleteRole(input: { id: $id })
		}
	`, variables, nil)
}

// GetRole returns a specific role
func (c *Client) GetRole(ctx context.Context, id database.ID) (*primitives.Role, error) {
	var resp struct {
		Role primitives.Role `json:"role"`
	}

	variables := make(map[string]interface{})
	variables["id"] = id

	err := c.Raw(ctx, `
		query Role($id: ID!) {
			role(id: $id) {
				id
				label
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Role, nil
}

// GetRoleByLabel returns a specific role by label
func (c *Client) GetRoleByLabel(ctx context.Context, label primitives.Label) (*primitives.Role, error) {
	var resp struct {
		Role primitives.Role `json:"roleByLabel"`
	}

	variables := make(map[string]interface{})
	variables["label"] = label

	err := c.Raw(ctx, `
		query RoleByLabel($label: Label!) {
			roleByLabel(label: $label) {
				id
				label
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Role, nil
}

// GetMembersRole returns the members of a role
func (c *Client) GetMembersRole(ctx context.Context, roleID database.ID) ([]primitives.Identity, error) {
	var resp struct {
		Identities []primitives.IdentityImpl `json:"roleMembers"`
	}

	variables := make(map[string]interface{})
	variables["role_id"] = roleID

	err := c.Raw(ctx, `
		query GetMembersRole($role_id: ID!) {
			roleMembers(role_id: $role_id) {
				id
				email
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return clientIdentitiesToPrimitive(resp.Identities)
}

func clientIdentitiesToPrimitive(identities []primitives.IdentityImpl) ([]primitives.Identity, error) {
	pIdentities := make([]primitives.Identity, len(identities))
	for i, identity := range identities {
		typ, err := identity.ID.Type()
		if err != nil {
			return nil, err
		}

		if typ == primitives.UserType {
			pIdentities[i] = &primitives.User{
				IdentityImpl: &identity,
			}
		} else if typ == primitives.ServicePrimitiveType {
			pIdentities[i] = &primitives.Service{
				IdentityImpl: &identity,
			}
		}
	}

	return pIdentities, nil
}

// AssignmentResponse is a type alias to easily decode the
// identity field to either a user or a service
type AssignmentResponse model.Assignment

// MarshalJSON marshaller impl for AssignmentResponse
func (a *AssignmentResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&model.Assignment{
		ID:        a.ID,
		Role:      a.Role,
		Identity:  a.Identity,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	})
}

// UnmarshalJSON is the marshaller implementation for AssignmentResponse
func (a *AssignmentResponse) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Identity *primitives.IdentityImpl `json:"identity"`
		Role     *primitives.Role         `json:"role"`
	}{}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	typ, err := aux.Identity.ID.Type()
	if err != nil {
		return err
	}

	if typ == primitives.UserType {
		a.Identity = &primitives.User{
			IdentityImpl: aux.Identity,
		}
	} else if typ == primitives.ServicePrimitiveType {
		a.Identity = &primitives.Service{
			IdentityImpl: aux.Identity,
		}
	}

	a.Role = aux.Role

	return nil
}

// AssignRole assigns a role to an identity
func (c *Client) AssignRole(ctx context.Context, identityID database.ID,
	roleID database.ID) (*model.Assignment, error) {
	var resp struct {
		Assignment AssignmentResponse `json:"assignRole"`
	}

	variables := make(map[string]interface{})
	variables["role_id"] = roleID
	variables["identity_id"] = identityID

	err := c.Raw(ctx, `
		mutation AssignRole($role_id: ID!, $identity_id: ID!) {
			assignRole(input: { role_id: $role_id, identity_id: $identity_id }) {
				role {
					id
					label
				}
				identity {
					id
					email
				}
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	assignment := model.Assignment(resp.Assignment)
	return &assignment, nil
}

// UnassignRole unassigns a role from an identity
func (c *Client) UnassignRole(ctx context.Context, identityID database.ID, roleID database.ID) error {
	variables := make(map[string]interface{})
	variables["role_id"] = roleID
	variables["identity_id"] = identityID

	return c.Raw(ctx, `
		mutation UnassignRole($role_id: ID!, $identity_id: ID!) {
			unassignRole(input: { role_id: $role_id, identity_id: $identity_id })
		}
	`, variables, nil)
}

// ListRoles returns all of the roles in the database
func (c *Client) ListRoles(ctx context.Context) ([]*primitives.Role, error) {
	var resp struct {
		Roles []*primitives.Role `json:"roles"`
	}

	err := c.Raw(ctx, `
		query Roles {
			roles {
				id
				label
			}
		}
	`, nil, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Roles, nil
}

// Source Routes

// SourceResponse is an alias of primitives.Source
// This is needed because the graphQL client cannot leverage the marshallers we have written
// for the URL properties of source (e.g. the Endpoint)
//
// We create a custom marshaller that encodes the endpoint as a string
type SourceResponse primitives.Source

// MarshalJSON is the marshaller implementation for SourceResponse
func (s *SourceResponse) MarshalJSON() ([]byte, error) {
	// We need another alias here as we are overwriting the Endpoint field of SourceResponse, which is a URL
	// If we embedded SourceResponse directly on the struct below, we would get an infinite loop trying to marshal
	// this object.  The type alias drops the Marshal & Unmarshal functions from "this" object.
	type SourceAlias SourceResponse
	return json.Marshal(&struct {
		Endpoint string `json:"endpoint"`
		*SourceAlias
	}{
		Endpoint:    s.Endpoint.String(),
		SourceAlias: (*SourceAlias)(s),
	})
}

// UnmarshalJSON is the marshaller implementation for SourceResponse
func (s *SourceResponse) UnmarshalJSON(data []byte) error {
	// See MarshalJSON for an explanation of this weird type alias
	type SourceAlias SourceResponse
	aux := &struct {
		Endpoint string `json:"endpoint"`
		*SourceAlias
	}{
		SourceAlias: (*SourceAlias)(s),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	e, err := url.Parse(aux.Endpoint)
	if err != nil {
		return err
	}

	s.Endpoint = *e
	return nil
}

// AddSource adds a new source to the database
func (c *Client) AddSource(ctx context.Context, label primitives.Label, credentials *url.URL) (*primitives.Source, error) {
	var resp struct {
		Source SourceResponse `json:"addSource"`
	}

	variables := make(map[string]interface{})
	variables["label"] = label
	variables["credentials"] = *credentials

	err := c.Raw(ctx, `
		mutation AddSource($label: Label!, $credentials: URL!) {
			  addSource(input: { label: $label, credentials: $credentials}) {
				id
				label
				endpoint
			  }
			}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	source := primitives.Source(resp.Source)
	return &source, nil
}

// ListSources returns all of the data sources in the database that you
func (c *Client) ListSources(ctx context.Context) ([]*primitives.Source, error) {
	var resp struct {
		Sources []SourceResponse `json:"sources"`
	}

	err := c.Raw(ctx, `
		query Sources {
				sources {
					id
					label
					endpoint
				}
			}
	`, nil, &resp)

	if err != nil {
		return nil, err
	}

	sources := make([]*primitives.Source, len(resp.Sources))
	for i := 0; i < len(sources); i++ {
		s := primitives.Source(resp.Sources[i])
		sources[i] = &s
	}

	return sources, nil
}

// GetSource returns a specific data source
func (c *Client) GetSource(ctx context.Context, id database.ID) (*primitives.Source, error) {
	var resp struct {
		Source SourceResponse `json:"source"`
	}

	variables := make(map[string]interface{})
	variables["id"] = id

	err := c.Raw(ctx, `
		query Sources($id: ID!) {
				source(id: $id) {
					id
					label
					endpoint
				}
			}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	source := primitives.Source(resp.Source)
	return &source, nil
}

// Setup calls the setup command to bootstrap cape
func (c *Client) Setup(ctx context.Context, user *primitives.User) (*primitives.User, error) {
	var resp struct {
		User primitives.User `json:"setup"`
	}

	variables := make(map[string]interface{})
	variables["name"] = user.Name
	variables["email"] = user.Email
	variables["public_key"] = user.Credentials.PublicKey
	variables["salt"] = user.Credentials.Salt

	err := c.Raw(ctx, `
		mutation CreateUser($name: Name!, $email: Email!, $public_key: Base64!, $salt: Base64!) {
			setup(input: { name: $name, email: $email, public_key: $public_key, salt: $salt, alg: "EDDSA"}) {
				id
				name
				email
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.User, nil
}

// CreateService creates a new service
func (c *Client) CreateService(ctx context.Context, service *primitives.Service) (*primitives.Service, error) {
	var resp struct {
		Service primitives.Service `json:"createService"`
	}

	variables := make(map[string]interface{})
	variables["email"] = service.Email
	variables["public_key"] = service.Credentials.PublicKey
	variables["salt"] = service.Credentials.Salt
	variables["type"] = service.Type

	err := c.Raw(ctx, `
		mutation CreateService($email: Email!, $type: ServiceType!, $public_key: Base64!, $salt: Base64!) {
			createService(input: { email: $email, type: $type, public_key: $public_key, salt: $salt, alg: "EDDSA"}) {
				id
				email
				type
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Service, nil
}

// DeleteService deletes a service
func (c *Client) DeleteService(ctx context.Context, id database.ID) error {
	variables := make(map[string]interface{})
	variables["id"] = id

	return c.Raw(ctx, `
		mutation DeleteService($id: ID!) {
			deleteService(input: { id: $id })
		}
	`, variables, nil)
}

// GetService returns a service by id
func (c *Client) GetService(ctx context.Context, id database.ID) (*primitives.Service, error) {
	var resp struct {
		Service primitives.Service `json:"service"`
	}

	variables := make(map[string]interface{})
	variables["id"] = id

	err := c.Raw(ctx, `
		query Service($id: ID!) {
			service(id: $id) {
				id
				email
				type
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Service, nil
}

// GetServiceByEmail returns a service by email
func (c *Client) GetServiceByEmail(ctx context.Context, email primitives.Email) (*primitives.Service, error) {
	var resp struct {
		Service primitives.Service `json:"serviceByEmail"`
	}

	variables := make(map[string]interface{})
	variables["email"] = email

	err := c.Raw(ctx, `
		query Service($email: Email!) {
			serviceByEmail(email: $email) {
				id
				email
				type
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Service, nil
}

// ListServices returns all services
func (c *Client) ListServices(ctx context.Context) ([]*primitives.Service, error) {
	var resp struct {
		Services []*primitives.Service `json:"services"`
	}

	err := c.Raw(ctx, `
		query Services {
			services {
				id
				email
				type
			}
		}
	`, nil, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Services, nil
}

// CreatePolicy creates a policy on the controller
func (c *Client) CreatePolicy(ctx context.Context, policy *primitives.Policy) (*primitives.Policy, error) {
	var resp struct {
		Policy primitives.Policy `json:"createPolicy"`
	}

	variables := make(map[string]interface{})
	variables["label"] = policy.Label

	err := c.Raw(ctx, `
		mutation CreatePolicy($label: Label!) {
			createPolicy(input: { label: $label }) {
				id
				label
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// DeletePolicy deletes a policy on the controller
func (c *Client) DeletePolicy(ctx context.Context, id database.ID) error {
	variables := make(map[string]interface{})
	variables["id"] = id

	return c.Raw(ctx, `
		mutation DeletePolicy($id: ID!) {
			deletePolicy(input: { id: $id })
		}
	`, variables, nil)
}

// GetPolicy returns a policy by id
func (c *Client) GetPolicy(ctx context.Context, id database.ID) (*primitives.Policy, error) {
	var resp struct {
		Policy primitives.Policy `json:"policy"`
	}

	variables := make(map[string]interface{})
	variables["id"] = id

	err := c.Raw(ctx, `
		query Policy($id: ID!) {
			policy(id: $id) {
				id
				label
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// ListPolicies returns all policies
func (c *Client) ListPolicies(ctx context.Context) ([]*primitives.Policy, error) {
	var resp struct {
		Policies []*primitives.Policy `json:"policies"`
	}

	err := c.Raw(ctx, `
		query Policies {
			policies {
				id
				label
			}
		}
	`, nil, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Policies, nil
}

// AttachPolicy attaches a policy to a role
func (c *Client) AttachPolicy(ctx context.Context, policyID database.ID,
	roleID database.ID) (*model.Attachment, error) {
	var resp struct {
		Attachment model.Attachment `json:"attachPolicy"`
	}

	variables := make(map[string]interface{})
	variables["role_id"] = roleID
	variables["policy_id"] = policyID

	err := c.Raw(ctx, `
		mutation AttachPolicy($role_id: ID!, $policy_id: ID!) {
			attachPolicy(input: { role_id: $role_id, policy_id: $policy_id }) {
				role {
					id
					label
				}
				policy {
					id
					label
				}
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Attachment, nil
}

// DetachPolicy unattaches a policy from a role
func (c *Client) DetachPolicy(ctx context.Context, policyID database.ID, roleID database.ID) error {
	variables := make(map[string]interface{})
	variables["role_id"] = roleID
	variables["policy_id"] = policyID

	return c.Raw(ctx, `
		mutation detachPolicy($role_id: ID!, $policy_id: ID!) {
			detachPolicy(input: { role_id: $role_id, policy_id: $policy_id })
		}
	`, variables, nil)
}

// GetRolePolicies returns all policies attached to a role
func (c *Client) GetRolePolicies(ctx context.Context, roleID database.ID) ([]*primitives.Policy, error) {
	var resp struct {
		Policies []*primitives.Policy `json:"rolePolicies"`
	}

	variables := make(map[string]interface{})
	variables["role_id"] = roleID

	err := c.Raw(ctx, `
		query RolePolicies($role_id: ID!) {
			rolePolicies(role_id: $role_id) {
				id
				label
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Policies, nil
}

// GetIdentityPolicies returns all policies related to an identity
func (c *Client) GetIdentityPolicies(ctx context.Context, identityID database.ID) ([]*primitives.Policy, error) {
	var resp struct {
		Policies []*primitives.Policy `json:"identityPolicies"`
	}

	variables := make(map[string]interface{})
	variables["identity_id"] = identityID

	err := c.Raw(ctx, `
		query IdentityPolicies($identity_id: ID!) {
			identityPolicies(identity_id: $identity_id) {
				id
				label
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Policies, nil
}
