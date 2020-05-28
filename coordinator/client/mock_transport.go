package coordinator

import (
	"context"
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/primitives"
	"github.com/manifoldco/go-base64"
	"reflect"
)

// MockTransport replaces the default transport on the client so we can return fake Responses to the CLI for testing
type MockTransport struct {
	Endpoint  *primitives.URL
	Responses []interface{}
	Counter   int
}

// Raw implements Raw on the Transport interface
func (m *MockTransport) Raw(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
	if len(m.Responses) == 0 {
		return nil
	}

	r := m.Responses[m.Counter]
	m.Counter++

	v := reflect.ValueOf(resp)
	v.Elem().Set(reflect.ValueOf(r))

	return nil
}

// Authenticated implements Authenticated on the Transport interface
func (m *MockTransport) Authenticated() bool {
	return true
}

// URL implements URL on the Transport interface
func (m *MockTransport) URL() *primitives.URL {
	return m.Endpoint
}

// EmailLogin implements EmailLogin on the Transport interface
func (m *MockTransport) EmailLogin(ctx context.Context, email primitives.Email, password primitives.Password) (*primitives.Session, error) {
	panic("Not Implemented")
}

// TokenLogin implements TokenLogin on the Transport interface
func (m *MockTransport) TokenLogin(ctx context.Context, token *auth.APIToken) (*primitives.Session, error) {
	panic("Not Implemented")
}

// Logout implements Logout on the Transport interface
func (m *MockTransport) Logout(ctx context.Context, authToken *base64.Value) error {
	panic("Not Implemented")
}
