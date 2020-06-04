package coordinator

import (
	"context"
	"net"
	"strings"

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
	return tokenLogin(ctx, c, apiToken)
}

// EmailLogin starts step 1 of the login flow using an email & password
func (c *HTTPTransport) EmailLogin(ctx context.Context, email primitives.Email, password primitives.Password) (*primitives.Session, error) {
	return emailLogin(ctx, c, email, password)
}

// Logout of the active session
func (c *HTTPTransport) Logout(ctx context.Context, authToken *base64.Value) error {
	return logout(ctx, c, authToken)
}

// Authenticated returns whether the client is authenticated or not.
// If the authToken is not nil then its authenticated!
func (c *HTTPTransport) Authenticated() bool {
	return c.authToken != nil
}

// SetToken enables a caller to set the auth token used by the transport
func (c *HTTPTransport) SetToken(value *base64.Value) {
	c.authToken = value
}

// Token enables a caller to retrieve the current auth token used by the
// transport
func (c *HTTPTransport) Token() *base64.Value {
	return c.authToken
}

func convertError(err error) error {
	return errors.New(errors.UnknownCause, strings.TrimPrefix(err.Error(), "graphql: "))
}
