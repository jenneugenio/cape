package primitives

import (
	"time"

	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/database/types"
	"github.com/manifoldco/go-base64"
)

// AuthCredentials represents the credentials the
// client will use to properly sign their challenge
type AuthCredentials struct {
	Salt *base64.Value
	Alg  CredentialsAlgType
}

// Session holds all the session data required to authenticate API
// calls with the server
type Session struct {
	*database.Primitive
	IdentityID database.ID   `json:"identity_id"`
	ExpiresAt  time.Time     `json:"expires_at"`
	Type       TokenType     `json:"type"`
	Token      *base64.Value `json:"token"`

	Credentials *AuthCredentials `json:"credentials"`
}

// GetType returns the type for this entity
func (s *Session) GetType() types.Type {
	return SessionType
}

// NewSession returns a new Session struct
func NewSession(identity Identity, expiresAt time.Time, typ TokenType,
	token *base64.Value) (*Session, error) {
	p, err := database.NewPrimitive(SessionType)
	if err != nil {
		return nil, err
	}

	session := &Session{
		Primitive:  p,
		IdentityID: identity.GetID(),
		ExpiresAt:  expiresAt,
		Type:       typ,
		Token:      token,

		Credentials: nil,
	}

	if typ == Login {
		creds := identity.GetCredentials()

		session.Credentials = &AuthCredentials{
			Salt: creds.Salt,
			Alg:  creds.Alg,
		}
	}

	id, err := database.DeriveID(session)
	if err != nil {
		return nil, err
	}
	session.ID = id

	return session, nil
}
