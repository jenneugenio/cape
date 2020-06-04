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
	SetToken(*base64.Value)
	Token() *base64.Value

	EmailLogin(ctx context.Context, email primitives.Email, password primitives.Password) (*primitives.Session, error)
	TokenLogin(ctx context.Context, apiToken *auth.APIToken) (*primitives.Session, error)

	Logout(ctx context.Context, authToken *base64.Value) error
}

func tokenLogin(ctx context.Context, transport ClientTransport, apiToken *auth.APIToken) (*primitives.Session, error) {
	var resp SessionResponse

	variables := map[string]interface{}{
		"token_id": apiToken.TokenID,
		"secret":   apiToken.Secret.Password(),
	}

	err := transport.Raw(ctx, `
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

	transport.SetToken(resp.Session.Token)
	return &resp.Session, nil
}

func emailLogin(ctx context.Context, transport ClientTransport, email primitives.Email, password primitives.Password) (*primitives.Session, error) {
	var resp SessionResponse

	variables := map[string]interface{}{
		"email":  email,
		"secret": password,
	}

	err := transport.Raw(ctx, `
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

	transport.SetToken(resp.Session.Token)
	return &resp.Session, nil
}

func logout(ctx context.Context, transport ClientTransport, authToken *base64.Value) error {
	token := transport.Token()
	if authToken != nil {
		token = authToken
	}

	variables := make(map[string]interface{})
	variables["token"] = token

	err := transport.Raw(ctx, `
		mutation DeleteSession($token: Base64) {
			deleteSession(input: { token: $token })
		}
	`, variables, nil)
	if token == transport.Token() {
		transport.SetToken(nil)
	}
	return err
}
