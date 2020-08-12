package coordinator

import (
	"context"
	"reflect"

	"github.com/manifoldco/go-base64"

	"encoding/json"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/models"
)

type MockRequest struct {
	Query     string
	Variables map[string]interface{}
	Body      interface{}
}

type MockResponse struct {
	Value interface{}
	Error error
}

// MockClientTransport replaces the default transport on the client so we can
// return fake Responses for unit testing
type MockClientTransport struct {
	Endpoint  *models.URL
	Requests  []*MockRequest
	Responses []*MockResponse
	Counter   int
	token     *base64.Value
}

func NewMockClientTransport(url *models.URL, responses []*MockResponse) (*MockClientTransport, error) {
	return &MockClientTransport{
		Endpoint:  url,
		Responses: responses,
		Counter:   0,
		token:     nil,
	}, nil
}

// Raw returns the appropriate response for the number request.
func (m *MockClientTransport) Raw(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
	m.Requests = append(m.Requests, &MockRequest{
		Query:     query,
		Variables: variables,
	})

	if len(m.Responses) == 0 {
		return nil
	}

	r := m.Responses[m.Counter]
	m.Counter++

	if r.Error != nil {
		return r.Error
	}

	if r.Value != nil {
		v := reflect.ValueOf(resp)
		v.Elem().Set(reflect.ValueOf(r.Value))
	}

	return nil
}

func (m *MockClientTransport) URL() *models.URL {
	return m.Endpoint
}

func (m *MockClientTransport) Authenticated() bool {
	return m.token != nil
}

func (m *MockClientTransport) SetToken(value *base64.Value) {
	m.token = value
}

func (m *MockClientTransport) Token() *base64.Value {
	return m.token
}

func (m *MockClientTransport) TokenLogin(ctx context.Context, apiToken *auth.APIToken) (*models.Session, error) {
	return tokenLogin(ctx, m, apiToken)
}

func (m *MockClientTransport) EmailLogin(ctx context.Context, email models.Email, password models.Password) (*models.Session, error) {
	return emailLogin(ctx, m, email, password)
}

func (m *MockClientTransport) Logout(ctx context.Context, authToken *base64.Value) error {
	return logout(ctx, m, authToken)
}

// Post does a raw http POST to the specified url
func (m *MockClientTransport) Post(url string, req interface{}) ([]byte, error) {
	m.Requests = append(m.Requests, &MockRequest{
		Body: req,
	})

	if len(m.Responses) == 0 {
		return nil, nil
	}

	r := m.Responses[m.Counter]
	m.Counter++

	if r.Error != nil {
		return nil, r.Error
	}

	by, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return by, nil
}
