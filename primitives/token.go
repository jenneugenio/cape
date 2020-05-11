package primitives

import (
	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/types"
)

type Token struct {
	*database.Primitive
	IdentityID database.ID        `json:"identity_id"`
	PublicKey  *base64.Value      `json:"public_key"`
	Salt       *base64.Value      `json:"salt"`
	Alg        CredentialsAlgType `json:"alg"`
}

func (tc *Token) GetType() types.Type {
	return TokenPrimitiveType
}

func (tc *Token) Validate() error {
	creds, err := tc.GetCredentials()
	if err != nil {
		return err
	}

	err = creds.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (tc *Token) GetIdentityID() database.ID {
	return tc.IdentityID
}

func (tc *Token) GetCredentials() (*Credentials, error) {
	return NewCredentials(tc.PublicKey, tc.Salt)
}

func NewToken(identityID database.ID, creds *Credentials) (*Token, error) {
	p, err := database.NewPrimitive(TokenPrimitiveType)
	if err != nil {
		return nil, err
	}

	tc := &Token{
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

func (tc *Token) GetEncryptable() bool {
	return false
}
