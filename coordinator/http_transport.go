package coordinator

import (
	"context"
	"net"

	"github.com/machinebox/graphql"
	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// HTTPTransport is a ClientTransport that interacts with the Coordinator via
// GraphQL over HTTP.
type HTTPTransport struct {
	client    *graphql.Client
	authToken *base64.Value
	url       *primitives.URL
}

// NewHTTPTransport returns a ClientTransport configured to make requests via
// GraphQL over HTTP
func NewHTTPTransport(coordinatorURL *primitives.URL, authToken *base64.Value) ClientTransport {
	if authToken != nil && len(authToken.String()) == 0 {
		authToken = nil
	}

	client := graphql.NewClient(coordinatorURL.String() + "/v1/query")
	return &HTTPTransport{
		client:    client,
		authToken: authToken,
		url:       coordinatorURL,
	}
}

// URL returns the underlying URL used by this Transport
func (c *HTTPTransport) URL() *primitives.URL {
	return c.url
}

// Raw wraps the NewRequest and does common req changes like adding authorization
// headers. It calls Run passing the object to be filled with the request data.
func (c *HTTPTransport) Raw(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
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

type SessionResponse struct {
	Session primitives.Session `json:"createSession"`
}

// TokenLogin enables a user or service to login using an APIToken
func (c *HTTPTransport) TokenLogin(ctx context.Context, apiToken *auth.APIToken) (*primitives.Session, error) {
	var resp SessionResponse

	variables := map[string]interface{}{
		"token_id": apiToken.TokenID,
		"secret":   apiToken.Secret.Password(),
	}

	err := c.Raw(ctx, `
		mutation CreateSession($token_id: ID, $secret: Password!) {
			createSession(input: { token_id: $token_id, secret: $secret }) {
				id
				identity_id
				expires_at
				token
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	c.authToken = resp.Session.Token
	return &resp.Session, nil
}

// EmailLogin starts step 1 of the login flow using an email & password
func (c *HTTPTransport) EmailLogin(ctx context.Context, email primitives.Email, password primitives.Password) (*primitives.Session, error) {
	var resp SessionResponse

	variables := map[string]interface{}{
		"email":  email,
		"secret": password,
	}

	err := c.Raw(ctx, `
		mutation CreateSession($email: Email, $secret: Password!) {
			createSession(input: { email: $email, secret: $secret }) {
				id
				identity_id
				expires_at
				token
			}
		}
	`, variables, &resp)
	if err != nil {
		return nil, err
	}

	c.authToken = resp.Session.Token
	return &resp.Session, nil
}

// Logout of the active session
func (c *HTTPTransport) Logout(ctx context.Context, authToken *base64.Value) error {
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
func (c *HTTPTransport) Authenticated() bool {
	return c.authToken != nil
}
