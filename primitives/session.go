package primitives

import (
	"time"

	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"

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

func (s *Session) Validate() error {
	if err := s.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidSessionCause, err)
	}

	if err := s.IdentityID.Validate(); err != nil {
		return errors.Wrap(InvalidSessionCause, err)
	}

	if time.Now().UTC().After(s.ExpiresAt) {
		return errors.New(InvalidSessionCause, "Session expires at must be after now")
	}

	if s.Token == nil {
		return errors.New(InvalidSessionCause, "Session token must not be nil")
	}

	return nil
}

// GetType returns the type for this entity
func (s *Session) GetType() types.Type {
	return SessionType
}

// NewSession returns a new Session struct
func NewSession(identity CredentialProvider, expiresAt time.Time, typ TokenType,
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
		creds, err := identity.GetCredentials()
		if err != nil {
			return nil, err
		}

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

	return session, session.Validate()
}
