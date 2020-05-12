package coordinator

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/manifoldco/go-base64"

	errors "github.com/capeprivacy/cape/partyerrors"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/graph/model"
	"github.com/capeprivacy/cape/primitives"

	"github.com/machinebox/graphql"
)

// NetworkCause occurs when the client cannot reach the server
var NetworkCause = errors.NewCause(errors.RequestTimeoutCategory, "network_error")

// UnrecognizedIdentityType occurs when the client encounters an identity type
// it doesn't recognize.
var UnrecognizedIdentityType = errors.NewCause(errors.BadRequestCategory, "unrecognized_identity")

type ClientTransport struct {
	client    *graphql.Client
	authToken *base64.Value
	url       *primitives.URL
}

func NewTransport(coordinatorURL *primitives.URL, authToken *base64.Value) Transport {
	if authToken != nil && len(authToken.String()) == 0 {
		authToken = nil
	}

	client := graphql.NewClient(coordinatorURL.String() + "/v1/query")
	//client.Log = func(s string) { fmt.Println("CLIENT TRANSPORT: ", s) }
	return &ClientTransport{
		client:    client,
		authToken: authToken,
		url:       coordinatorURL,
	}
}

// createTokenLoginSession runs a CreateTokenLoginSession mutation that creates a
// login session and returns it
func (c *ClientTransport) createTokenLoginSession(ctx context.Context, token *auth.APIToken) (*primitives.Session, error) {
	var resp struct {
		Session primitives.Session `json:"createLoginSession"`
	}

	variables := make(map[string]interface{})
	variables["token_id"] = token.TokenID

	err := c.Raw(ctx, `
		mutation CreateLoginSession($token_id: ID!) {
			createLoginSession(input: { token_id: $token_id }) {
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

// createEmailLoginSession runs a CreateLoginSession mutation that creates a
// login session and returns it
func (c *ClientTransport) createEmailLoginSession(ctx context.Context, email primitives.Email) (*primitives.Session, error) {
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
func (c *ClientTransport) createAuthSession(ctx context.Context, sig *base64.Value) (*primitives.Session, error) {
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

func (c *ClientTransport) URL() *primitives.URL {
	return c.url
}

// Raw wraps the NewRequest and does common req changes like adding authorization
// headers. It calls Run passing the object to be filled with the request data.
func (c *ClientTransport) Raw(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
	req := graphql.NewRequest(query)

	for key, val := range variables {
		req.Var(key, val)
	}

	if c.authToken != nil {
		req.Header.Add("Authorization", "Bearer "+c.authToken.String())
	}

	err := c.client.Run(ctx, req, resp)
	if err != nil {
		if nerr, ok := err.(net.Error); ok {
			return errors.New(NetworkCause, "Could not contact coordinator: %s", nerr.Error())
		}
		return convertError(err)
	}

	return nil
}

// TokenLogin starts step 1 of the login flow using an API token
func (c *ClientTransport) TokenLogin(ctx context.Context, apiToken *auth.APIToken) (*primitives.Session, error) {
	session, err := c.createTokenLoginSession(ctx, apiToken)
	if err != nil {
		return nil, err
	}

	c.authToken = session.Token

	creds, err := auth.NewCredentials(apiToken.Secret, session.Credentials.Salt)
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

// EmailLogin starts step 1 of the login flow using an email & password
func (c *ClientTransport) EmailLogin(ctx context.Context, email primitives.Email, password auth.Secret) (*primitives.Session, error) {
	session, err := c.createEmailLoginSession(ctx, email)
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

// Logout of the active session
func (c *ClientTransport) Logout(ctx context.Context, authToken *base64.Value) error {
	var token *base64.Value
	if authToken != nil {
		token = authToken
	}

	variables := make(map[string]interface{})
	variables["token"] = token

	return c.Raw(ctx, `
		mutation DeleteSession($token: Base64) {
			deleteSession(input: { token: $token })
		}
	`, variables, nil)
}

// Authenticated returns whether the client is authenticated or not.
// If the authToken is not nil then its authenticated!
func (c *ClientTransport) Authenticated() bool {
	return c.authToken != nil
}

// Client is a wrapper around the graphql client that
// connects to the coordinator and sends queries
type Client struct {
	transport Transport
}

// NewClient returns a new client that connects to the given
// the configured transport
func NewClient(transport Transport) *Client {
	return &Client{
		transport: transport,
	}
}

func convertError(err error) error {
	return errors.New(errors.UnknownCause, strings.TrimPrefix(err.Error(), "graphql: "))
}

type MeResponse struct {
	Identity *primitives.IdentityImpl `json:"me"`
}

// Me returns the identity of the current authenticated session
func (c *Client) Me(ctx context.Context) (primitives.Identity, error) {
	var resp MeResponse

	err := c.transport.Raw(ctx, `
		query Me {
			me {
				id
				email
			}
		}
	`, nil, &resp)
	if err != nil {
		return nil, err
	}

	return identityImplToIdentity(resp.Identity)
}

// UserResponse is a primitive User with an extra Roles field that maps to the
// GraphQL type.
type UserResponse struct {
	*primitives.User
	Roles []*primitives.Role `json:"roles"`
}

// GetUser returns a user and it's roles!
func (c *Client) GetUser(ctx context.Context, id database.ID) (*UserResponse, error) {
	var resp struct {
		User UserResponse `json:"user"`
	}

	variables := make(map[string]interface{})
	variables["id"] = id.String()

	err := c.transport.Raw(ctx, `
		query User($id: ID!) {
			user(id: $id) {
				id
				name
				email
				roles {
					id
					label
				}
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.User, nil
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

	err := c.transport.Raw(ctx, `
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

func (c *Client) Authenticated() bool {
	return c.transport.Authenticated()
}

// EmailLogin calls the CreateLoginSession and CreateAuthSession mutations
func (c *Client) EmailLogin(ctx context.Context, email primitives.Email, password auth.Secret) (*primitives.Session, error) {
	return c.transport.EmailLogin(ctx, email, password)
}

func (c *Client) TokenLogin(ctx context.Context, token *auth.APIToken) (*primitives.Session, error) {
	return c.transport.TokenLogin(ctx, token)
}

// Logout calls the deleteSession mutation
func (c *Client) Logout(ctx context.Context, authToken *base64.Value) error {
	return c.transport.Logout(ctx, authToken)
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

	err := c.transport.Raw(ctx, `
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

	return c.transport.Raw(ctx, `
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

	err := c.transport.Raw(ctx, `
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

	err := c.transport.Raw(ctx, `
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
		Identities []*primitives.IdentityImpl `json:"roleMembers"`
	}

	variables := make(map[string]interface{})
	variables["role_id"] = roleID

	err := c.transport.Raw(ctx, `
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

func clientIdentitiesToPrimitive(identities []*primitives.IdentityImpl) ([]primitives.Identity, error) {
	out := make([]primitives.Identity, len(identities))
	for i, identity := range identities {
		result, err := identityImplToIdentity(identity)
		if err != nil {
			return nil, err
		}

		out[i] = result
	}

	return out, nil
}

func identityImplToIdentity(identity *primitives.IdentityImpl) (primitives.Identity, error) {
	typ, err := identity.ID.Type()
	if err != nil {
		return nil, err
	}

	switch typ {
	case primitives.UserType:
		return &primitives.User{IdentityImpl: identity}, nil
	case primitives.ServicePrimitiveType:
		return &primitives.Service{IdentityImpl: identity}, nil
	default:
		return nil, errors.New(UnrecognizedIdentityType, "Unknown Type: %s", typ.String())
	}
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

	err := c.transport.Raw(ctx, `
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

	return c.transport.Raw(ctx, `
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

	err := c.transport.Raw(ctx, `
		query Roles {
			roles {
				id
				label
				system
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
type SourceResponse struct {
	*primitives.Source
	Service *primitives.Service `json:"service"`
}

// MarshalJSON is the marshaller implementation for SourceResponse
//func (s *SourceResponse) MarshalJSON() ([]byte, error) {
//	// We need another alias here as we are overwriting the Endpoint field of SourceResponse, which is a URL
//	// If we embedded SourceResponse directly on the struct below, we would get an infinite loop trying to marshal
//	// this object.  The type alias drops the Marshal & Unmarshal functions from "this" object.
//	type SourceAlias SourceResponse
//	return json.Marshal(&struct {
//		Endpoint    string `json:"endpoint"`
//		Credentials string `json:"credentials"`
//		*SourceAlias
//	}{
//		Endpoint:    s.Endpoint.String(),
//		Credentials: s.Credentials.String(),
//		SourceAlias: (*SourceAlias)(s),
//	})
//}
//
//// UnmarshalJSON is the marshaller implementation for SourceResponse
//func (s *SourceResponse) UnmarshalJSON(data []byte) error {
//	// See MarshalJSON for an explanation of this weird type alias
//	type SourceAlias SourceResponse
//	aux := &struct {
//		Endpoint    string `json:"endpoint"`
//		Credentials string `json:"credentials"`
//		*SourceAlias
//	}{
//		SourceAlias: (*SourceAlias)(s),
//	}
//
//	err := json.Unmarshal(data, &aux)
//	if err != nil {
//		return err
//	}
//
//	e, err := primitives.NewDBURL(aux.Endpoint)
//	if err != nil {
//		return err
//	}
//
//	if aux.Credentials != "" {
//		c, err := primitives.NewDBURL(aux.Credentials)
//		if err != nil {
//			return err
//		}
//
//		s.Credentials = c
//	}
//
//	s.Endpoint = e
//	return nil
//}

// AddSource adds a new source to the database
func (c *Client) AddSource(ctx context.Context, label primitives.Label,
	credentials *primitives.DBURL, serviceID *database.ID) (*SourceResponse, error) {
	var resp struct {
		Source SourceResponse `json:"addSource"`
	}

	variables := make(map[string]interface{})
	variables["label"] = label
	variables["credentials"] = credentials.String()
	variables["service_id"] = serviceID

	err := c.transport.Raw(ctx, `
		mutation AddSource($label: Label!, $credentials: DBURL!, $service_id: ID) {
			  addSource(input: { label: $label, credentials: $credentials, service_id: $service_id}) {
				id
				label
				type
				credentials
				endpoint
				service {
					id
					email
				}
			  }
			}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Source, nil
}

type ListSourcesResponse struct {
	Sources []*SourceResponse `json:"sources"`
}

// ListSources returns all of the data sources in the database that you
func (c *Client) ListSources(ctx context.Context) ([]*SourceResponse, error) {
	var resp ListSourcesResponse
	err := c.transport.Raw(ctx, `
		query Sources {
				sources {
					id
					label
					credentials
					type
					endpoint
					service {
						id
						email
						endpoint
					}
				}
			}
	`, nil, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Sources, nil
}

func (c *Client) RemoveSource(ctx context.Context, label primitives.Label) error {
	variables := make(map[string]interface{})
	variables["label"] = label

	// We only care if this errors
	var resp interface{}
	return c.transport.Raw(ctx, `
		mutation RemoveSource($label: Label!) {
			  removeSource(input: { label: $label })
			}
	`, variables, &resp)
}

// GetSource returns a specific data source
func (c *Client) GetSource(ctx context.Context, id database.ID) (*SourceResponse, error) {
	var resp struct {
		Source SourceResponse `json:"source"`
	}

	variables := make(map[string]interface{})
	variables["id"] = id

	err := c.transport.Raw(ctx, `
		query Sources($id: ID!) {
				source(id: $id) {
					id
					label
					type
					credentials
					endpoint
					service {
						id
						email
					}
				}
			}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Source, nil
}

// GetSourceByLabel returns a specific data source given its label
func (c *Client) GetSourceByLabel(ctx context.Context, label primitives.Label) (*SourceResponse, error) {
	var resp struct {
		Source SourceResponse `json:"sourceByLabel"`
	}

	variables := make(map[string]interface{})
	variables["label"] = label

	fmt.Println("no transport??", c.transport)
	err := c.transport.Raw(ctx, `
		query Sources($label: Label!) {
				sourceByLabel(label: $label) {
					id
					label
					type
					credentials
					endpoint
					service {
						id
						email
					}
				}
			}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Source, nil
}

// Setup calls the setup command to bootstrap cape
func (c *Client) Setup(ctx context.Context, user *primitives.User) (*primitives.User, error) {
	var resp struct {
		User *primitives.User `json:"setup"`
	}

	in := &model.NewUserRequest{
		Name:      user.Name,
		Email:     user.Email,
		PublicKey: *user.Credentials.PublicKey,
		Salt:      *user.Credentials.Salt,
		Alg:       primitives.EDDSA,
	}

	variables := make(map[string]interface{})
	variables["input"] = in

	err := c.transport.Raw(ctx, `
		mutation CreateUser($input: NewUserRequest!) {
			setup(input: $input) {
				id
				name
				email
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

// ServiceResponse is a primitive Service with an extra Roles field
type ServiceResponse struct {
	*primitives.Service
	Roles []*primitives.Role `json:"roles"`
}

// CreateService creates a new service
func (c *Client) CreateService(ctx context.Context, service *primitives.Service) (*primitives.Service, error) {
	var resp struct {
		Service ServiceResponse `json:"createService"`
	}

	variables := make(map[string]interface{})
	variables["email"] = service.Email
	variables["type"] = service.Type

	variables["endpoint"] = nil
	if service.Endpoint != nil {
		variables["endpoint"] = service.Endpoint.String()
	}

	err := c.transport.Raw(ctx, `
		mutation CreateService($email: Email!, $type: ServiceType!, $endpoint: URL) {
			createService(input: { email: $email, type: $type, endpoint: $endpoint }) {
				id
				email
				type
				endpoint
			}
		}
	`, variables, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Service.Service, nil
}

// DeleteService deletes a service
func (c *Client) DeleteService(ctx context.Context, id database.ID) error {
	variables := make(map[string]interface{})
	variables["id"] = id

	return c.transport.Raw(ctx, `
		mutation DeleteService($id: ID!) {
			deleteService(input: { id: $id })
		}
	`, variables, nil)
}

// GetService returns a service by id
func (c *Client) GetService(ctx context.Context, id database.ID) (*primitives.Service, error) {
	var resp struct {
		Service ServiceResponse `json:"service"`
	}

	variables := make(map[string]interface{})
	variables["id"] = id

	err := c.transport.Raw(ctx, `
		query Service($id: ID!) {
			service(id: $id) {
				id
				email
				type
				endpoint
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Service.Service, nil
}

// GetServiceByEmail returns a service by email
func (c *Client) GetServiceByEmail(ctx context.Context, email primitives.Email) (*primitives.Service, error) {
	var resp struct {
		Service ServiceResponse `json:"serviceByEmail"`
	}

	variables := make(map[string]interface{})
	variables["email"] = email

	err := c.transport.Raw(ctx, `
		query Service($email: Email!) {
			serviceByEmail(email: $email) {
				id
				email
				type
				endpoint
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Service.Service, nil
}

// ListServices returns all services
func (c *Client) ListServices(ctx context.Context) ([]*ServiceResponse, error) {
	var resp struct {
		Services []*ServiceResponse `json:"services"`
	}

	err := c.transport.Raw(ctx, `
		query Services {
			services {
				id
				email
				type
				endpoint
				roles {
					id
					label
				}
			}
		}
	`, nil, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Services, nil
}

// CreatePolicy creates a policy on the coordinator
func (c *Client) CreatePolicy(ctx context.Context, policy *primitives.Policy) (*primitives.Policy, error) {
	var resp struct {
		Policy primitives.Policy `json:"createPolicy"`
	}

	variables := make(map[string]interface{})
	variables["label"] = policy.Label
	variables["spec"] = policy.Spec

	err := c.transport.Raw(ctx, `
		mutation CreatePolicy($label: Label!, $spec: PolicySpec!) {
			createPolicy(input: { label: $label, spec: $spec }) {
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

// DeletePolicy deletes a policy on the coordinator
func (c *Client) DeletePolicy(ctx context.Context, id database.ID) error {
	variables := make(map[string]interface{})
	variables["id"] = id

	return c.transport.Raw(ctx, `
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

	err := c.transport.Raw(ctx, `
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

// GetPolicyByLabel returns a specific policy by label
func (c *Client) GetPolicyByLabel(ctx context.Context, label primitives.Label) (*primitives.Policy, error) {
	var resp struct {
		Policy primitives.Policy `json:"policyByLabel"`
	}

	variables := make(map[string]interface{})
	variables["label"] = label

	err := c.transport.Raw(ctx, `
		query PolicyByLabel($label: Label!) {
			policyByLabel(label: $label) {
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

	err := c.transport.Raw(ctx, `
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

	err := c.transport.Raw(ctx, `
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

	return c.transport.Raw(ctx, `
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

	err := c.transport.Raw(ctx, `
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

	err := c.transport.Raw(ctx, `
		query IdentityPolicies($identity_id: ID!) {
			identityPolicies(identity_id: $identity_id) {
				id
				label
				spec
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Policies, nil
}

// GetIdentities returns all identities for the given emails
func (c *Client) GetIdentities(ctx context.Context, emails []primitives.Email) ([]primitives.Identity, error) {
	var resp struct {
		Identities []*primitives.IdentityImpl `json:"identities"`
	}

	variables := make(map[string]interface{})
	variables["emails"] = emails

	err := c.transport.Raw(ctx, `
		query IdentityPolicies($emails: [Email!]) {
			identities(emails: $emails) {
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

type CreateTokenResponse struct {
	Token primitives.Token `json:"createToken"`
}

// CreateToken creates a new API token for the provided identity. You can pass nil and it will return a token for you
func (c *Client) CreateToken(ctx context.Context, identity primitives.Identity) (*auth.APIToken, error) {
	// If the user provides no identity, we will make a token for the current session user
	if identity == nil {
		i, err := c.Me(ctx)
		if err != nil {
			return nil, err
		}

		identity = i
	}

	secret, err := auth.RandomSecret()
	if err != nil {
		return nil, err
	}

	creds, err := auth.NewCredentials(secret, nil)
	if err != nil {
		return nil, err
	}

	pCreds, err := creds.Package()
	if err != nil {
		return nil, err
	}

	tokenCredentials, err := primitives.NewToken(identity.GetID(), pCreds)
	if err != nil {
		return nil, err
	}

	variables := make(map[string]interface{})
	variables["identity_id"] = tokenCredentials.IdentityID
	variables["public_key"] = tokenCredentials.PublicKey
	variables["salt"] = tokenCredentials.Salt
	variables["alg"] = tokenCredentials.Alg

	var resp CreateTokenResponse
	err = c.transport.Raw(ctx, `
		mutation CreateToken($identity_id: ID!, $public_key: Base64!, $salt: Base64!, $alg: CredentialsAlgType!) {
			createToken(input: { identity_id: $identity_id, public_key: $public_key, salt: $salt, alg: $alg }) {
				id
			}
        }
    `, variables, &resp)
	if err != nil {
		return nil, err
	}

	token, err := auth.NewAPIToken(secret, resp.Token.ID, c.transport.URL())
	return token, err
}

type ListTokensResponse struct {
	IDs []database.ID `json:"tokens"`
}

// ListTokens lists all of the auth tokens for the provided identity
func (c *Client) ListTokens(ctx context.Context, identity primitives.Identity) ([]database.ID, error) {
	// If the user provides no identity, we will make a token for the current session user
	if identity == nil {
		i, err := c.Me(ctx)
		if err != nil {
			return nil, err
		}

		identity = i
	}

	var resp ListTokensResponse

	variables := make(map[string]interface{})
	variables["identity_id"] = identity.GetID()

	err := c.transport.Raw(ctx, `
		query Tokens($identity_id: ID!) {
			tokens(identity_id: $identity_id)
		}
    `, variables, &resp)

	if err != nil {
		return nil, err
	}

	return resp.IDs, nil
}

// RemoveTokens removes the provided token from the database
func (c *Client) RemoveToken(ctx context.Context, tokenID database.ID) error {
	variables := make(map[string]interface{})
	variables["id"] = tokenID

	return c.transport.Raw(ctx, `
		mutation RemoveToken($id: ID!) {
			removeToken(id: $id)
		}
    `, variables, nil)
}
