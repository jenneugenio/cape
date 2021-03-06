package coordinator

import (
	"context"
	"encoding/json"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/models"
)

// ClientTransport is an interface that describes how a coordinator client should communicate with a coordinator
type ClientTransport interface {
	Raw(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error

	// Post does a raw http POST to the specified url
	Post(url string, req interface{}) ([]byte, error)

	Authenticated() bool
	URL() *models.URL
	SetToken(*base64.Value)
	Token() *base64.Value

	EmailLogin(ctx context.Context, email models.Email, password models.Password) (*models.Session, error)
	TokenLogin(ctx context.Context, apiToken *auth.APIToken) (*models.Session, error)

	Logout(ctx context.Context, authToken *base64.Value) error
}

func tokenLogin(ctx context.Context, transport ClientTransport, apiToken *auth.APIToken) (*models.Session, error) {
	req := LoginRequest{
		TokenID: &apiToken.TokenID,
		Secret:  apiToken.Secret.Password(),
	}

	body, err := transport.Post(transport.URL().String()+"/v1/login", req)
	if err != nil {
		return nil, err
	}

	session := &models.Session{}
	err = json.Unmarshal(body, session)
	if err != nil {
		return nil, err
	}

	transport.SetToken(session.Token)
	return session, nil
}

func emailLogin(ctx context.Context, transport ClientTransport, email models.Email, password models.Password) (*models.Session, error) {
	req := LoginRequest{
		Email:  &email,
		Secret: password,
	}

	body, err := transport.Post(transport.URL().String()+"/v1/login", req)
	if err != nil {
		return nil, err
	}

	session := &models.Session{}
	err = json.Unmarshal(body, session)
	if err != nil {
		return nil, err
	}
	transport.SetToken(session.Token)

	return session, nil
}

func logout(ctx context.Context, transport ClientTransport, authToken *base64.Value) error {
	token := transport.Token()
	if authToken != nil {
		token = authToken
	}

	req := LogoutRequest{
		Token: token,
	}

	_, err := transport.Post(transport.URL().String()+"/v1/logout", req)
	if err != nil {
		return err
	}

	if token == transport.Token() {
		transport.SetToken(nil)
	}

	return err
}
