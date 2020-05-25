package primitives

import (
	"context"
	"encoding/json"

	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

type Token struct {
	*database.Primitive
	IdentityID database.ID `json:"identity_id"`

	// We never want to send Credentials over the wire!
	Credentials *Credentials `json:"credentials" gqlgen:"-"`
}

type encryptedToken struct {
	*Token
	Credentials *base64.Value `json:"credentials"`
}

func (tc *Token) GetType() types.Type {
	return TokenPrimitiveType
}

func (tc *Token) Validate() error {
	if err := tc.Primitive.Validate(); err != nil {
		return err
	}

	if tc.Credentials == nil {
		return errors.New(InvalidTokenCause, "Credentials must be non-nil")
	}

	if err := tc.IdentityID.Validate(); err != nil {
		return err
	}

	t, err := tc.IdentityID.Type()
	if err != nil {
		return err
	}

	if t != UserType && t != ServicePrimitiveType {
		return errors.New(InvalidTokenCause, "IdentityID must be a user or service")
	}

	return tc.Credentials.Validate()
}

func (tc *Token) GetIdentityID() database.ID {
	return tc.IdentityID
}

func (tc *Token) GetCredentials() (*Credentials, error) {
	return tc.Credentials, nil
}

func NewToken(identityID database.ID, creds *Credentials) (*Token, error) {
	p, err := database.NewPrimitive(TokenPrimitiveType)
	if err != nil {
		return nil, err
	}

	tc := &Token{
		Primitive:   p,
		IdentityID:  identityID,
		Credentials: creds,
	}

	id, err := database.DeriveID(tc)
	if err != nil {
		return nil, err
	}

	tc.ID = id
	return tc, tc.Validate()
}

func (tc *Token) Encrypt(ctx context.Context, codec crypto.EncryptionCodec) ([]byte, error) {
	creds, err := json.Marshal(tc.Credentials)
	if err != nil {
		return nil, err
	}

	data, err := codec.Encrypt(ctx, base64.New(creds))
	if err != nil {
		return nil, err
	}

	return json.Marshal(encryptedToken{
		Token:       tc,
		Credentials: data,
	})
}

func (tc *Token) Decrypt(ctx context.Context, codec crypto.EncryptionCodec, data []byte) error {
	in := &encryptedToken{}
	err := json.Unmarshal(data, in)
	if err != nil {
		return err
	}

	unencrypted, err := codec.Decrypt(ctx, in.Credentials)
	if err != nil {
		return err
	}

	creds := &Credentials{}
	err = json.Unmarshal([]byte(*unencrypted), creds)
	if err != nil {
		return err
	}

	tc.Primitive = in.Primitive
	tc.IdentityID = in.IdentityID

	tc.Credentials = creds
	return nil
}

func (tc *Token) GetEncryptable() bool {
	return true
}

// GenerateToken returns an instantiated token for use in unit testing.
//
// This function _should only ever_ be used inside of a test.
func GenerateToken(identity Identity) (Password, *Token, error) {
	password, err := GeneratePassword()
	if err != nil {
		return EmptyPassword, nil, err
	}

	c, err := GenerateCredentials()
	if err != nil {
		return EmptyPassword, nil, err
	}

	token, err := NewToken(identity.GetID(), c)
	return password, token, err
}
