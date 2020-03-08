package primitives

import (
	"github.com/dropoutlabs/privacyai/primitives/types"
)

// Token for an authorized entity (user or service)
type Token struct {
	*Primitive
	IdentityID ID `json:"identity_id"`
}

// GetType returns the type for this entity
func (t *Token) GetType() types.Type {
	return TokenType
}

// NewToken returns an immutable token struct
func NewToken(identityID ID) (*Token, error) {
	p, err := newPrimitive(TokenType)
	if err != nil {
		return nil, err
	}

	// TODO: Figure out what fields should be set on this struct
	t := &Token{
		Primitive:  p,
		IdentityID: identityID,
	}

	ID, err := DeriveID(t)
	if err != nil {
		return nil, err
	}

	t.ID = ID
	return t, nil
}
