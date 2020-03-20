package controller

import (
	"context"
	"fmt"
	"net/url"

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
func (c *Client) Raw(ctx context.Context, query string, resp interface{}) error {
	req := graphql.NewRequest(query)

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
	`, email), &resp)

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
	`, sig.String()), &resp)

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
	`, user.Name, user.Email, user.Credentials.PublicKey.String(), user.Credentials.Salt.String()), &resp)

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

	err := c.Raw(ctx, fmt.Sprintf(`
		mutation DeleteSession {
		  deleteSession(input: { token: "%s" })
		}
	`, token), nil)

	if err != nil {
		return err
	}

	return nil
}
