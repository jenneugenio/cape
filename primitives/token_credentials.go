package primitives

import (
	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/database/types"
	"github.com/manifoldco/go-base64"
)

type TokenCredentials struct {
	*database.Primitive
	IdentityID database.ID        `json:"identity_id"`
	PublicKey  *base64.Value      `json:"public_key"`
	Salt       *base64.Value      `json:"salt"`
	Alg        CredentialsAlgType `json:"alg"`
}

func (tc *TokenCredentials) GetType() types.Type {
	return TokenPrimitiveType
}

func (tc *TokenCredentials) Validate() error {
	return nil
}

func (tc *TokenCredentials) GetCredentials() (*Credentials, error) {
	return NewCredentials(tc.PublicKey, tc.Salt)
}

func NewTokenCredentials(identityID database.ID, creds *Credentials) (*TokenCredentials, error) {
	p, err := database.NewPrimitive(TokenPrimitiveType)
	if err != nil {
		return nil, err
	}

	tc := &TokenCredentials{
		Primitive:  p,
		IdentityID: identityID,
		PublicKey:  creds.PublicKey,
		Salt:       creds.Salt,
		Alg:        creds.Alg,
	}

	id, err := database.DeriveID(tc)
	if err != nil {
		return nil, err
	}

	tc.ID = id
	return tc, tc.Validate()
}
