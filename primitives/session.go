package primitives

import (
	"context"
	"encoding/json"
	"time"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

// Session holds all the session data required to authenticate API
// calls with the server
type Session struct {
	*database.Primitive
	IdentityID database.ID   `json:"identity_id"`
	OwnerID    database.ID   `json:"owner_id"`
	ExpiresAt  time.Time     `json:"expires_at"`
	Token      *base64.Value `json:"token"`
}

type encryptedSession struct {
	*Session
	Token *base64.Value `json:"token"`
}

func (s *Session) Validate() error {
	if err := s.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidSessionCause, err)
	}

	if err := s.IdentityID.Validate(); err != nil {
		return errors.Wrap(InvalidSessionCause, err)
	}

	identityTypes := []types.Type{
		UserType,
		ServicePrimitiveType,
	}
	if !s.IdentityID.OneOf(identityTypes) {
		return errors.New(InvalidSessionCause, "Identity ID is not a user or service")
	}

	if err := s.OwnerID.Validate(); err != nil {
		return errors.Wrap(InvalidSessionCause, err)
	}

	if !s.OwnerID.OneOf([]types.Type{UserType, TokenPrimitiveType}) {
		return errors.New(InvalidSessionCause, "Owner ID is not a user or token")
	}

	return nil
}

// GetType returns the type for this entity
func (s *Session) GetType() types.Type {
	return SessionType
}

// NewSession returns a new Session struct
func NewSession(identity CredentialProvider) (*Session, error) {
	p, err := database.NewPrimitive(SessionType)
	if err != nil {
		return nil, err
	}

	session := &Session{
		Primitive:  p,
		IdentityID: identity.GetIdentityID(),
		OwnerID:    identity.GetID(),
	}

	id, err := database.DeriveID(session)
	if err != nil {
		return nil, err
	}
	session.ID = id

	return session, session.Validate()
}

// Encrypt implements the Encryptable interface
func (s *Session) Encrypt(ctx context.Context, codec crypto.EncryptionCodec) ([]byte, error) {
	data, err := codec.Encrypt(ctx, s.Token)
	if err != nil {
		return nil, err
	}

	return json.Marshal(encryptedSession{
		Session: s,
		Token:   data,
	})
}

// Decrypt implements the Encryptable interface
func (s *Session) Decrypt(ctx context.Context, codec crypto.EncryptionCodec, data []byte) error {
	in := &encryptedSession{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	unencrypted, err := codec.Decrypt(ctx, in.Token)
	if err != nil {
		return err
	}

	s.Primitive = in.Primitive
	s.IdentityID = in.IdentityID
	s.OwnerID = in.OwnerID
	s.ExpiresAt = in.ExpiresAt

	s.Token = unencrypted
	return nil
}

func (s *Session) SetToken(token *base64.Value, expiresAt time.Time) {
	s.Token = token
	s.ExpiresAt = expiresAt
}

func (s *Session) GetEncryptable() bool {
	return true
}
