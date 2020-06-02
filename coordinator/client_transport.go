package coordinator

import (
	"context"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/primitives"
)

// ClientTransport is an interface that describes how a coordinator client should communicate with a coordinator
type ClientTransport interface {
	Raw(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error
	Authenticated() bool
	URL() *primitives.URL

	EmailLogin(ctx context.Context, email primitives.Email, password primitives.Password) (*primitives.Session, error)
	TokenLogin(ctx context.Context, apiToken *auth.APIToken) (*primitives.Session, error)

	Logout(ctx context.Context, authToken *base64.Value) error
}
