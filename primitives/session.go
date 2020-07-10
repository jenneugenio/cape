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
	UserID    string        `json:"user_id"`
	OwnerID   string        `json:"owner_id"`
	ExpiresAt time.Time     `json:"expires_at"`
	Token     *base64.Value `json:"token"`
}

type encryptedSession struct {
	*Session
	Token *base64.Value `json:"token"`
}

func (s *Session) Validate() error {
	if err := s.Primitive.Validate(); err != nil {
		return errors.Wrap(InvalidSessionCause, err)
	}

	if s.UserID == "" {
		return errors.New(InvalidSessionCause, "IdentityID is empty")
	}

	if s.OwnerID == "" {
		return errors.New(InvalidSessionCause, "OwnerID is empty")
	}

	return nil
}

// GetType returns the type for this entity
func (s *Session) GetType() types.Type {
	return SessionType
}

// NewSession returns a new Session struct
func NewSession(cp CredentialProvider) (*Session, error) {
	p, err := database.NewPrimitive(SessionType)
	if err != nil {
		return nil, err
	}

	session := &Session{
		Primitive: p,
		UserID:    cp.GetUserID(),
		OwnerID:   cp.GetStringID(),
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
	s.UserID = in.UserID
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
