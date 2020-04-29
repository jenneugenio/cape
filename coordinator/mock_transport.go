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
	Responses []interface{}
	Counter   int
}

// Raw implements Raw on the Transport interface
func (m MockTransport) Raw(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
	r := m.Responses[m.Counter]
	m.Counter++

	v := reflect.ValueOf(resp)
	v.Elem().Set(reflect.ValueOf(r))

	return nil
}

// Authenticated implements Authenticated on the Transport interface
func (m MockTransport) Authenticated() bool {
	return true
}

// EmailLogin implements EmailLogin on the Transport interface
func (m MockTransport) Login(ctx context.Context, email primitives.Email, password auth.Secret) (*primitives.Session, error) {
	panic("Not Implemented")
}

// Logout implements Logout on the Transport interface
func (m MockTransport) Logout(ctx context.Context, authToken *base64.Value) error {
	panic("Not Implemented")
}
