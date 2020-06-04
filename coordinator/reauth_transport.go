package coordinator

import (
	"context"
	"strings"
	"sync"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// ReAuthTransport is a ClientTransport that wraps an HTTPTransport to
// reauthenticate when a session token expires.
//
// ReAuthTransport only supports authentication through an API Token.
type ReAuthTransport struct {
	transport ClientTransport
	apiToken  *auth.APIToken

	// ReAuthTransport uses a mutex to avoid multiple calls from attempting to
	// authenticate at the same time.
	mutex *sync.Mutex
}

// NewReAuthTransport returns a ReAuthTransport
func NewReAuthTransport(coordinatorURL *primitives.URL, apiToken *auth.APIToken) (ClientTransport, error) {
	if coordinatorURL == nil {
		return nil, errors.New(InvalidArgumentCause, "Missing coordinator url")
	}

	if err := coordinatorURL.Validate(); err != nil {
		return nil, err
	}

	if apiToken == nil {
		return nil, errors.New(InvalidArgumentCause, "Missing apiToken")
	}

	if err := apiToken.Validate(); err != nil {
		return nil, err
	}

	return &ReAuthTransport{
		transport: NewHTTPTransport(coordinatorURL, nil),
		apiToken:  apiToken,
		mutex:     &sync.Mutex{},
	}, nil
}

// Raw enables a Client to perform a raw GraphQL request against a server
func (r *ReAuthTransport) Raw(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
	// If we're not authenticated, then let's attempt to authenticate!
	if !r.Authenticated() {
		err := r.attempt(ctx)
		if err != nil {
			return err
		}
	}

	// If we get an authentication error back then we attempt to login. If that
	// fails, we bubble the error back up.
	err := r.transport.Raw(ctx, query, variables, resp)
	if isAuthenticationError(err) {
		err = r.attempt(ctx)
		if err != nil {
			return err
		}

		return r.transport.Raw(ctx, query, variables, resp)
	}

	return err
}

func (r *ReAuthTransport) attempt(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, err := r.transport.TokenLogin(ctx, r.apiToken)
	return err
}

// EmailLogin is not implemented by ReAuthTransport. Only authentication
// through an APIToken is supported.
func (r *ReAuthTransport) EmailLogin(ctx context.Context, email primitives.Email, password primitives.Password) (*primitives.Session, error) {
	return nil, errors.ErrNotImplemented
}

// TokenLogin is not implemented by ReAuthTransport. You must supply the
// APIToken at instantiation.
func (r *ReAuthTransport) TokenLogin(ctx context.Context, apiToken *auth.APIToken) (*primitives.Session, error) {
	return nil, errors.ErrNotImplemented
}

// Authenticated returns whether or not the client is currently authenticated
func (r *ReAuthTransport) Authenticated() bool {
	return r.transport.Authenticated()
}

// URL returns the underlying URL used by the Transport
func (r *ReAuthTransport) URL() *primitives.URL {
	return r.transport.URL()
}

// Logout attempts to log out of the current session
func (r *ReAuthTransport) Logout(ctx context.Context, authToken *base64.Value) error {
	return r.transport.Logout(ctx, authToken)
}

// Token enables a caller to access the current token from the transport
func (r *ReAuthTransport) Token() *base64.Value {
	return r.transport.Token()
}

// SetToken enables a caller to set the auth token used by the transport
func (r *ReAuthTransport) SetToken(value *base64.Value) {
	r.transport.SetToken(value)
}

func isAuthenticationError(err error) bool {
	e, ok := err.(*errors.Error)
	if !ok {
		return false
	}

	// XXX: We are not propagating the "cause" (e.g. 401 Unauthorized) through
	// our GQL client. Therefore, we have to rely _only_ on the message.
	for _, msg := range e.Messages {
		if strings.Contains(msg, auth.ErrAuthentication.Messages[0]) {
			return true
		}
	}

	return false
}
