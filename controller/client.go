package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/dropoutlabs/cape/database"

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
func (c *Client) createLoginSession(ctx context.Context, email string) (*primitives.Session, error) {
	var resp struct {
		Session primitives.Session `json:"createLoginSession"`
	}

	err := c.Raw(ctx, fmt.Sprintf(`
		mutation CreateLoginSession{
			createLoginSession(input: { email: "%s" }) {
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
	`, email), nil, &resp)

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

	err := c.Raw(ctx, fmt.Sprintf(`
		mutation CreateAuthSession{
			createAuthSession(input: { signature: "%s" }) {
				id
				identity_id
				expires_at
				type
				token
			}
		}
	`, sig.String()), nil, &resp)

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

	err := c.Raw(ctx, fmt.Sprintf(`
		mutation CreateUser {
			createUser(input: { name: "%s", email: "%s", public_key: "%s", salt: "%s", alg: "EDDSA"}) {
				id
				name
				email
			}
		}
	`, user.Name, user.Email, user.Credentials.PublicKey.String(), user.Credentials.Salt.String()), nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.User, nil
}

// Login calls the CreateLoginSession and CreateAuthSession
// mutations
func (c *Client) Login(ctx context.Context, email string, password []byte) (*primitives.Session, error) {
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
	var token *base64.Value
	if authToken == nil {
		token = c.authToken
	}

	return c.Raw(ctx, fmt.Sprintf(`
		mutation DeleteSession {
			deleteSession(input: { token: "%s" })
		}
	`, token), nil, nil)
}

// Role Routes

// CreateRole creates a new role with a label
func (c *Client) CreateRole(ctx context.Context, label string, identityIDs []database.ID) (*primitives.Role, error) {
	var resp struct {
		Role primitives.Role `json:"createRole"`
	}

	variables := make(map[string]interface{})
	variables["ids"] = identityIDs

	err := c.Raw(ctx, fmt.Sprintf(`
		mutation CreateRole($ids: [ID!]) {
			createRole(input: { label: "%s", identity_ids: $ids }) {
				id
				label
			}
		}
	`, label), variables, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Role, nil
}

// DeleteRole deletes a role with the given id
func (c *Client) DeleteRole(ctx context.Context, id database.ID) error {
	return c.Raw(ctx, fmt.Sprintf(`
		mutation DeleteRole {
			deleteRole(input: { id: "%s" })
		}
	`, id.String()), nil, nil)
}

// GetRole returns a specific role
func (c *Client) GetRole(ctx context.Context, id database.ID) (*primitives.Role, error) {
	var resp struct {
		Role primitives.Role `json:"role"`
	}

	err := c.Raw(ctx, fmt.Sprintf(`
		query Role {
			role(id: "%s") {
				id
				label
			}
		}
	`, id.String()), nil, &resp)
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

	err := c.Raw(ctx, fmt.Sprintf(`
		query GetMembersRole {
			roleMembers(role_id: "%s") {
				id
				email
			}
		}
	`, roleID.String()), nil, &resp)
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
		} else if typ == primitives.ServiceType {
			pIdentities[i] = &primitives.Service{
				IdentityImpl: &identity,
			}
		}
	}

	return pIdentities, nil
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
func (c *Client) AddSource(ctx context.Context, label string, credentials *url.URL) (*primitives.Source, error) {
	var resp struct {
		Source SourceResponse `json:"addSource"`
	}

	err := c.Raw(ctx, fmt.Sprintf(`
		mutation AddSource {
			  addSource(input: { label: "%s", credentials: "%s"}) {
				id
				label
				endpoint
			  }
			}
	`, label, credentials.String()), nil, &resp)
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

	err := c.Raw(ctx, fmt.Sprintf(`
		query Sources {
				source(id: "%s") {
					id
					label
					endpoint
				}
			}
	`, id.String()), nil, &resp)
	if err != nil {
		return nil, err
	}

	source := primitives.Source(resp.Source)
	return &source, nil
}
